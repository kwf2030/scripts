package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

var dir string

func main() {
	flag.Parse()
	dir = flag.Arg(0)
	if dir == "" {
		dir = "files"
	}
	f, e := os.Open(dir)
	if e != nil {
		createDir()
	} else if s, _ := f.Stat(); !s.IsDir() {
		createDir()
	}
	fmt.Printf("Start serving directory [%s]\n", dir)
	serve()
}

func createDir() {
	e := os.Mkdir(dir, 0666)
	if e != nil {
		panic(e)
	}
}

func serve() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.Write([]byte("Requires GET"))
			return
		}
		w.Write([]byte(index_tpl))
	})
	http.HandleFunc("/upload/", uploadHandler)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(dir))))
	http.ListenAndServe(":8082", nil)
}

var uploadHandler = func(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Write([]byte("Requires POST"))
		return
	}
	r.ParseMultipartForm(32 << 20)
	src, h, e := r.FormFile("file")
	if e != nil {
		fmt.Println(e)
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dir+"/"+h.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	fmt.Printf("[%s]\n", dst.Name())
	if err != nil {
		fmt.Println(e)
		return
	}
	defer dst.Close()
	io.Copy(dst, src)
}

const index_tpl = `
<!DOCTYPE html>
<html lang="zh">
<head>
  <meta charset="UTF-8">
  <title>Simple File Server</title>
  <style type="text/css">
    pre a {
      font-size: 1.2em;
      line-height: 150%;
    }
    #progress {
      color: #9F9F9F;
      font-size: .8em;
      font-style: italic;
    }
  </style>
</head>
<body onload="getFileList()">
<div>
  <form action="/upload/" enctype="multipart/form-data" method="POST">
    <input type="file" id="file"/>
    <p></p>
    <div>
    <input type="button" onclick="uploadFile()" value="上传"/>
    <label id="progress"></label>
    </div>
  </form>
</div>
<hr />
<p></p>
<div id="files"></div>
<script type="text/javascript">
  var progress = document.getElementById('progress')

  function uploadFile() {
    var file = document.getElementById('file').files[0];
    if (!file) {
      alert("未选择任何文件");
      return;
    }
    var form = new FormData();
    form.append("file", file);
    var xhr = new XMLHttpRequest();
    xhr.upload.addEventListener("progress", function(evt) {
      if (evt.lengthComputable) {
        var r = Math.round(evt.loaded * 100 / evt.total);
        progress.innerHTML = r.toString() + '%';
      }
    }, false);

    xhr.addEventListener("load", function(evt) {
      progress.innerHTML = 'Success';
    }, false);

    xhr.addEventListener("error", function(evt) {
      progress.innerHTML = 'Failed';
    }, false);

    xhr.addEventListener("abort", function(evt) {
      progress.innerHTML = 'Canceled';
    }, false);

    xhr.open("POST", "/upload/");
    xhr.send(form);
  }

  function getFileList() {
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function() {
      if (xhr.readyState == 4 && xhr.status == 200) {
        var result = xhr.responseText.replace(/href="/g, "download href=\"files/")
        document.getElementById("files").innerHTML = result;
      }
    }
    xhr.open("GET", "/files/");
    xhr.send();
  }
</script>
</body>
</html>
`
