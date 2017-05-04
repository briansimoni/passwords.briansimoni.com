# CMSC 413 Final Project

## Simoni Passwords

Simoni Passwords is a web application designed for secure storage of secret information, specifically passwords. At the time of this writing, the project is deployed to https://passwords.briansimoni.com/


### Design Overview and architecture Diagram

![architecture](https://github.com/briansimoni/passwords.briansimoni.com/blob/master/architecture.png)

### Operating System

The application is cross platform ready, but the current production build is deployed on an Ubuntu 16 Linux instance on Amazon Web services Elastic Compute Cloud. Amazon is configured to block all traffic except for ports 443, 80, and 22. (http later is force upgraded to https). The application files themselves are protected with basic Linux file system permissions. The master key has a permissions of 400 and is owned by the “Jenkins” user.

Database

The database choice is MySQL and follows a fairly basic security setup. Remote connections are not allowed, and users have only the necessary permissions. Application code utilizes prepared statements to prevent SQL Injection

Web Server

Apache2 is used as the public facing web server. It will handle all of the connections to ports 80 and 443, and then act as a reverse proxy to the application running on port 8888. Apache also enforces https by automatically redirecting connections to port 80 to port 443.

Server Side Application Code

The application itself is written in Go (https://golang.org/) - an upcoming language created by Google. Go is a compiled, strongly typed language ideal for web, distributed, and concurrent applications. It is additionally garbage collected and safe from buffer overflow attacks since it performs bounds checking. When a user signs up for the application, the bcrypt library is utilized to hash the password with a random salt and add additional computational overhead to prevent viable computation of things like rainbow tables. When users try to authenticate, it will simply hash the provided password and compare the hash to the one stored in the database.
Encryption is done with AES-256. Both the user’s application names and secrets are encrypted. At no point in time are secrets ever stored in plaintext on persistent storage. Sessions are also encrypted with the same algorithm.

### Continuous Integration

To speed up the development process, the Jenkins continuous integration server is utilized for automatic deployments. Whenever a developer makes a code commit to the Github master branch, a new build is triggered. Jenkins will automatically fetch the latest code, compile, and run the application.


### End user security and experience

The web application provides a sleek design that is responsive to both desktop and mobile phone screen sizes. Since it is on the public internet, it is always at a convenient location for users on any device. On the client side, there is basic form validation to prevent the user from submitting bad requests. It also enforces a strong password, and provides helpful messages when appropriate. All connections made are https to prevent eavesdropping attacks.



### Local Application Setup

Requirements
-	Golang installed
-	MySQL installed locally
o	MySQL database created with the name SimoniPasswords and appropriate tables
-	Add a master password file

Once you have met the requirements, you can compile with either go build or go install



### Recommendation for Grading

Since this project uses non-standard technologies and it may be difficult to compile and run locally, you can always access the source code at 

https://github.com/briansimoni/passwords.briansimoni.com

and the application is 

https://passwords.briansimoni.com/


You can email any questions to me at simonibc@vcu.edu

