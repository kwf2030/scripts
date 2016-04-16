# scripts

## publish-to-maven.gradle

A simple gradle script that can publish your local artifacts to maven repositories.

Follow these steps:

1. Create a 'libs' directory that contains all artifacts to publish
2. Run `gradle` or `gradle conf` to generate 'conf.json'
3. Configure 'conf.json', fill in groupId/artifactId/version
4. Run `gradle publishToMavenLocal` or `publishToMavenRepository` to publish your artifacts

Done.

## netty-build.gradle & netty-settings.gradle

Gradle build script for netty

## fs.go

A simple file server that support download/upload.

Use `fs -h` to see the usage.
