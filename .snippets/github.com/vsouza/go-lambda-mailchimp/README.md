# Golang AWS Lambda Mailchimp subscriber

[![License][license-image]][license-url]

> Simple AWS lambda written in golang to subscribe users on Mailchimp list

To deploy and test you can use [Apex](https://github.com/apex) :)

## Install

At the first, you need to install [Apex](https://github.com/apex), it's very simple:

`curl https://raw.githubusercontent.com/apex/apex/master/install.sh | sh`

After configure your AWS Credentials, then run:

`apex deploy`

## Run Locally

create a `event.json` file like this:

```
{"email_address": "email@vsouza.com", "status": "subscribed", "merge_fields": {"FNAME": "Vinicius Souza"}}
```

save on project root folder.

Then run:

`apex invoke newsletter-subscribe < event.json`

## Release History

* 0.0.1
    * Subscribe users with first name on mailchimp 3.0 API

## Meta

Vinicius Souza – [@iamvsouza](https://twitter.com/iamvsouza) – hi@vsouza.com

Distributed under the MIT license. See [License](http://vsouza.mit-license.org/)

[https://github.com/vsouza](https://github.com/vsouza/)

[license-image]: https://img.shields.io/badge/License-MIT-blue.svg
[license-url]: LICENSE
