import java.nio.file.Files

import de.svenjacobs.loremipsum.LoremIpsum
import devsearch.concat.{ Main, Utils }
import org.scalatest.{ Matchers, WordSpecLike }

import scala.sys.process._

/**
 * Created by dengels on 03/05/15.
 */
class SystemSpec extends WordSpecLike with Matchers {

  "The system" should {

    "Have the correct size for small text files" in {

      val loremIpsum = new LoremIpsum()
      val rootDir = Files.createTempDirectory("rootDir")
      val langDir = Files.createDirectory(rootDir.resolve("lang"))
      val ownerDir = Files.createDirectory(langDir.resolve("owner"))
      val files = for (i <- 1 to 10) yield {
        val repo = Files.createDirectory(ownerDir.resolve(s"repo_$i"))
        val file = Files.createFile(repo.resolve("file"))
        val bytes = loremIpsum.getParagraphs(i).getBytes("UTF-8")
        Files.write(file, bytes)
      }

      val totalSize = files.map { file => Files.size(file) }.sum

      val resDir = Files.createTempDirectory("resDir")

      Main.main(Array("-j", "1") ++ Array(rootDir, resDir).map(_.toAbsolutePath.toString))

      val blobs = resDir.toFile.listFiles()

      blobs.size should equal(1)

      val resSize = Files.size(blobs(0).toPath)

      resSize should be > totalSize

    }

    "Be able to read tarball repos" in {

      val loremIpsum = new LoremIpsum()
      val rootDir = Files.createTempDirectory("rootDir")
      val langDir = Files.createDirectory(rootDir.resolve("lang"))
      val ownerDir = Files.createDirectory(langDir.resolve("owner"))
      val repo = Files.createDirectory(ownerDir.resolve(s"repo"))
      val files = for (i <- 1 to 10) yield {
        val file = Files.createFile(repo.resolve(s"file$i"))
        val bytes = loremIpsum.getParagraphs(i).getBytes("UTF-8")
        Files.write(file, bytes)
      }

      val totalSize = files.map { file => Files.size(file) }.sum

      /* Create tar */
      val command = Seq("tar", "-cf", "repo.tar", repo.getFileName.toString)
      val rv = Process(command, cwd = Some(ownerDir.toFile)).!

      rv should be(0)

      println(s"Created tarball : ${ownerDir.resolve("repo.tar")}")

      for (file <- Utils.listFilesRec(repo)) Files.delete(file)
      Files.delete(repo)

      val resDir = Files.createTempDirectory("resDir")

      Main.main(Array("-j", "1") ++ Array(rootDir, resDir).map(_.toAbsolutePath.toString))

      val blobs = resDir.toFile.listFiles()

      blobs.size should equal(1)

      val resSize = Files.size(blobs(0).toPath)

      resSize should be > totalSize

    }

  }
}
