# goFire

Agile development management dashboard with emphasis on administrative-level analysis.

View burndown and overview charts for all developers, teams, projects and milestones throughout your agile development process.

"It's all about the burndown baby."


# Screenshots

## People

![alt tag](https://raw.githubusercontent.com/gophergala/goFire/master/galainfo/people.png)

## Sprints

![alt tag](https://raw.githubusercontent.com/gophergala/goFire/master/galainfo/sprints.png)

![alt tag](https://raw.githubusercontent.com/gophergala/goFire/master/galainfo/sprintdetail.png)

## Tasks

![alt tag](https://raw.githubusercontent.com/gophergala/goFire/master/galainfo/taskview.png)

![alt tag](https://raw.githubusercontent.com/gophergala/goFire/master/galainfo/taskcomment.png)


# Limitations

Currently a lot of pages don't actually post back / get data from a database. These need to be completed. Burndown charts and styling on the dashboard are not yet done. Analytics page, auth, and account pages are not finished / started.

# To Get Started

You need the following env variables setup on your machine

1. CONNECTIONSTRING = "Your mysql connection string"
2. PORT = "Port to listen on for this app when started"
3. COOKIESTORE = "Cookie Secret"
4. DB who's schema matches the struct setup in the model directory