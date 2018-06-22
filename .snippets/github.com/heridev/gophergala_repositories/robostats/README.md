# Robostats

Robotstats is a platform to gather and display statistics for robots and other
connected devices.

~~Robostats can also be used as a remote control to send signals to connected
devices of the same kind, when a device receives this signal it can alter its
behaviour and provide feedback.~~

![general](https://cloud.githubusercontent.com/assets/385670/5893726/9ff493ac-a4b4-11e4-8a94-d769a9d0efea.png)

You may see a development preview at
[dev.robostats.io](http://dev.robostats.io). Use `user@example.com` / `pass` to
enter, we didn't have any security precautions, so please don't hack us :-).

## The API

![api](https://cloud.githubusercontent.com/assets/385670/5886015/140d0e7c-a348-11e4-9f29-2e1599f98994.png)

The [Robostats API](https://github.com/gophergala/robostats/tree/master/go/src/robostats/api)
provides tools for creating categories for your devices (or robot *classes*).
Each class has its own API key, connected devices can use this API key to
register themselves (as a robot *instance*).

A robot instance must create a session before adding events, a session is like
a bucket to store actual device logs.

Different devices may produce different events and event values, that's why we
use MongoDB as permanent storage, its ability to save schemaless documents can
be used to store many different kinds of logs.

## Client package

There is also a
[client](https://github.com/gophergala/robostats/tree/master/go/src/robostats/client)
package for Robostats that can be used as an example to communicate with the
API.

## Panel for analytics

![frontend](https://cloud.githubusercontent.com/assets/385670/5886013/14073830-a348-11e4-961a-fdc9bde0b151.png)

The analytics panel provides an overview of all registered devices, organized
by it's class. It allows users to register new device classes and get the
generated api key for each class.

![web-1](https://cloud.githubusercontent.com/assets/1133/5893820/7531d81a-a4b8-11e4-9316-7126455d99c5.png)

Selecting a device class users can drill down to a specific device, it's
session (any period of activity) and finally, to a graph that shows all the
events generated from that drone in a particular sesion.

![web-2](https://cloud.githubusercontent.com/assets/1133/5893822/75379016-a4b8-11e4-864b-3d931eb8fcc0.png)

![web-3](https://cloud.githubusercontent.com/assets/1133/5893821/7533b572-a4b8-11e4-8885-5682af246575.png)

## Remote control

The remove control as not finished on time, we left it here as a reference for
future implementation.

![remote-control](https://cloud.githubusercontent.com/assets/385670/5886016/181bf8ca-a348-11e4-89bc-4463931da63d.png)

The [remote
control](https://github.com/gophergala/robostats/tree/master/android) is an
application for Robostats that can be used to send signals to connected
devices. Robostats only acts as a distribution hub for these commands.

The WIP implementation of the remote controller was an Android app.

![android-1](https://cloud.githubusercontent.com/assets/385670/5894130/685688ca-a4c2-11e4-9b4c-cf01e70a2961.png)

You could use the app to log in into your account.

![android-2](https://cloud.githubusercontent.com/assets/385670/5894129/6853c86a-a4c2-11e4-9d6f-639cc282b442.png)

Using the remote control, you could choose a device class and send commands to
all connected instances. These commands could vary, but we we're playing with a
standard game-like pad.

![android-3](https://cloud.githubusercontent.com/assets/385670/5894128/6852c7c6-a4c2-11e4-84a9-194a52eae3b4.png)

## Internal documentation

If you're interested on hacking the platform, you may want to read these
resources:

* [API endpoints](https://github.com/gophergala/robostats/tree/master/doc)

## License

> Copyright 2015. The robostats authors.
>
> Licensed under the Apache License, Version 2.0 (the "License");
> you may not use this file except in compliance with the License.
> You may obtain a copy of the License at
>
>     http://www.apache.org/licenses/LICENSE-2.0
>
> Unless required by applicable law or agreed to in writing, software
> distributed under the License is distributed on an "AS IS" BASIS,
> WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
> See the License for the specific language governing permissions and
> limitations under the License.
