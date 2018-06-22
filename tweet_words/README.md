Sentiment Analysis
==================

* To run this application, you need to have mongoDB installed

* You also need to create a twitter application at https://dev.twitter.com and copy the information to a new file called app.properties as below.

* Create file app.properties with credentials before executing with  following keys.

CONSUMER_KEY=
CONSUMER_SECRET=
ACCESS_TOKEN=
ACCESS_TOKEN_SECRET=
MONGO=mongodb://

* You don't need to specify the database name in the MONGO url. The application uses the database name "test" for now.

* Go to the classifier_train directory and create the classifier.gob file as per the instructions in the README file. Copy the classifier.gob to the executable directory

* After you have app.properties and classifier.gob ready. You can execute the application by running.


Copy /bower_components/ and main.go outside tweetwords folder and run
> go run main.go
