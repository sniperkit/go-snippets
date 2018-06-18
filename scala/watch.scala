#!/usr/bin/env scala

import java.nio.file._
import scala.collection.JavaConversions._
import scala.sys.process._

val file = Paths.get(args(0))
val cmd = args(1)
val watcher = FileSystems.getDefault.newWatchService

file.register(
  watcher, 
  StandardWatchEventKinds.ENTRY_CREATE,
  StandardWatchEventKinds.ENTRY_MODIFY,
  StandardWatchEventKinds.ENTRY_DELETE
)

def exec = cmd run true

@scala.annotation.tailrec
def watch(proc: Process): Unit = {
  val key = watcher.take
  val events = key.pollEvents

  val newProc = 
    if (!events.isEmpty) {
      proc.destroy()
      exec
    } else proc

  if (key.reset) watch(newProc)
  else println("aborted")
}

watch(exec)
