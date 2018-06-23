#!/bin/sh

show_usuage() 
{
    echo "Usuage: gnew.sh ProjectName"
}

if [ $# -ne 1 ]; then
    show_usuage;
    exit -1;
fi

project=$1;
package=`echo $project | tr '[:upper:]' '[:lower:]'`

mkdir $1
cd $1

# Create Project class
mkdir -p src/main/java/$package
touch src/main/java/$package/$1.java
echo "package $package;"            > src/main/java/$package/$1.java
echo ""                             >> src/main/java/$package/$1.java
echo "public class $project {"      >> src/main/java/$package/$1.java
echo "    public $project() {"      >> src/main/java/$package/$1.java
echo "    }"                        >> src/main/java/$package/$1.java
echo "}"                            >> src/main/java/$package/$1.java

# Create Project Main class
touch src/main/java/$package/$1Main.java
echo "package $package;"            > src/main/java/$package/$1Main.java
echo ""                             >> src/main/java/$package/$1Main.java
echo "public class $1Main {"        >> src/main/java/$package/$1Main.java
echo "    public static void main(String[] args) {"     >> src/main/java/$package/$1Main.java
echo "        System.out.println(\"Hello $project!\"); "  >> src/main/java/$package/$1Main.java
echo "    }"                        >> src/main/java/$package/$1Main.java
echo "}"                            >> src/main/java/$package/$1Main.java

# Create Project Test class
mkdir -p src/test/java/$package
touch src/test/java/$package/$1Test.java
echo "package $package;"            > src/test/java/$package/$1Test.java
echo ""                             >> src/test/java/$package/$1Test.java
echo "import org.junit.Before;"     >> src/test/java/$package/$1Test.java
echo "import org.junit.Test;"       >> src/test/java/$package/$1Test.java
echo "import static org.junit.Assert.assertEquals;"     >> src/test/java/$package/$1Test.java
echo ""                             >> src/test/java/$package/$1Test.java
echo "public class $1Test {"        >> src/test/java/$package/$1Test.java
echo "    @Test"                    >> src/test/java/$package/$1Test.java
echo "    public void TestTest() {" >> src/test/java/$package/$1Test.java
echo "        assertEquals(\"Hello World\", \"Hello world\");" >> src/test/java/$package/$1Test.java
echo "    }"                        >> src/test/java/$package/$1Test.java
echo "}"                            >> src/test/java/$package/$1Test.java

# Add Gradle configuration file
touch build.gradle
echo "apply plugin : 'java'          "> build.gradle
echo "apply plugin : 'maven'         ">> build.gradle
echo "apply plugin : 'application'   ">> build.gradle
echo "                               ">> build.gradle
echo "repositories {                 ">> build.gradle
echo "    mavenCentral()             ">> build.gradle
echo "}                              ">> build.gradle
echo "                               ">> build.gradle
echo "dependencies {                 ">> build.gradle
echo "    compile (                  ">> build.gradle
echo "    'log4j:log4j:1.2.17'       ">> build.gradle
echo "    )                          ">> build.gradle
echo "    testCompile (              ">> build.gradle
echo "    'junit:junit:4.+'          ">> build.gradle
echo "    )                          ">> build.gradle
echo "}                              ">> build.gradle
echo "                               ">> build.gradle
echo "jar {                          ">> build.gradle
echo "    baseName = 'kkkmmu'        ">> build.gradle
echo "    version = '0.1.0'          ">> build.gradle
echo "}                              ">> build.gradle
echo "                               ">> build.gradle
echo "task wrapper(type: Wrapper) {  ">> build.gradle
echo "    gradleVersion = '3.0'      ">> build.gradle
echo "}                              ">> build.gradle
echo "                               ">> build.gradle
echo "mainClassName = '$package.$1Main' ">> build.gradle

gradle wrapper
