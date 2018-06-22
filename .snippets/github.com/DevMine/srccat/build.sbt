import scalariform.formatter.preferences._

name := """srccat"""

shellPrompt := { state => "[\033[36m" + name.value + "\033[0m] $ " }

version := "1.0.0"

scalaVersion := "2.11.6"

scalacOptions ++= Seq("-unchecked", "-deprecation", "-feature")

libraryDependencies ++= Seq(
  "com.typesafe.akka" %% "akka-actor" % "2.3.9",
  "com.typesafe.akka" %% "akka-testkit" % "2.3.9" % "test",
  "org.scalatest" %% "scalatest" % "2.2.4" % "test",
  "com.github.scopt" %% "scopt" % "3.3.0",
  "org.apache.tika" % "tika-core" % "1.7",
  "de.sven-jacobs" % "loremipsum" % "1.0",
  "org.apache.commons" % "commons-compress" % "1.9",
  "org.apache.commons" % "commons-io" % "1.3.2"
)

resolvers += Resolver.sonatypeRepo("public")

target in Compile in doc := (baseDirectory.value / "api")

scalariformSettings

ScalariformKeys.preferences := ScalariformKeys.preferences.value
    .setPreference(AlignParameters, true)
    .setPreference(AlignSingleLineCaseStatements, false)
    .setPreference(CompactControlReadability, false)
    .setPreference(CompactStringConcatenation, false)
    .setPreference(DoubleIndentClassDeclaration, true)
    .setPreference(FormatXml, true)
    .setPreference(IndentLocalDefs, false)
    .setPreference(IndentPackageBlocks, true)
    .setPreference(IndentSpaces, 2)
    .setPreference(IndentWithTabs, false)
    .setPreference(MultilineScaladocCommentsStartOnFirstLine, false)
    .setPreference(PreserveDanglingCloseParenthesis, false)
    .setPreference(RewriteArrowSymbols, false)
    .setPreference(SpaceBeforeColon, false)
    .setPreference(SpaceInsideBrackets, false)
    .setPreference(SpaceInsideParentheses, false)
    .setPreference(SpacesWithinPatternBinders, true)

wartremoverErrors ++= Warts.allBut(Wart.DefaultArguments, Wart.Nothing, Wart.FinalCaseClass, Wart.NoNeedForMonad, Wart.Any, Wart.Throw)
wartremoverWarnings ++= Warts.allBut(Wart.FinalCaseClass, Wart.NoNeedForMonad, Wart.Any, Wart.Throw)
