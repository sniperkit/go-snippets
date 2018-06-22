package devmine.srccat.actors

import devmine.srccat.actors.Worker.{ BlobRequest, Begin }
import akka.actor.ActorSystem
import akka.testkit.{ TestKit, ImplicitSender }
import org.scalatest.WordSpecLike
import org.scalatest.Matchers
import org.scalatest.BeforeAndAfterAll

class WorkerSpec(_system: ActorSystem) extends TestKit(_system) with ImplicitSender
    with WordSpecLike with Matchers with BeforeAndAfterAll {

  def this() = this(ActorSystem("MySpec"))

  override def afterAll {
    TestKit.shutdownActorSystem(system)
  }

  "A Worker actor" must {
    "send a Blob request on Begin" in {
      val pingActor = system.actorOf(Worker.props(self))
      pingActor ! Begin
      val msg = expectMsg(BlobRequest)
    }
  }

}
