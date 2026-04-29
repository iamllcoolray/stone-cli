plugins {
    id 'java'
    id 'application'
}

group '{{.PackageName}}'
version '1.0'

repositories {
    mavenCentral()
    maven { url 'https://maven.pkg.jetbrains.space/litiengine/p/maven/releases' }
}

dependencies {
    implementation 'de.gurkenlabs:litiengine:{{.LitiVersion}}'
}

application {
    mainClass = '{{.PackageName}}.Main'
}
