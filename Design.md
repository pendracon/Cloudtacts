# Cloudtacts
## Summary
A simple contacts database manager for the cloud.

## Scope
Cloudtacts provides a simple, multiuser contacts database maintained in the
cloud and accessible through RESTful APIs. Cloudtacts is initially conceived
as a learning exercise for cloud native development and also as an aid to
preparing for Google Cloud Engineer and Cloud Developer certifications.
Consequently, initial deployment is targeted to Google Cloud, utilizing
Google Cloud services and features. Ports to AWS, Azure, and IBM clouds might
be pursued at a later time for similar respective purposes.

## Requirements
TODO

## Architecture
Cloudtacts includes the following components:
1. A SQL database containing user information.
2. A NoSQL database for storing user contacts as JSON documents.
3. An object store for storing profile images associated with users and
their contacts.
4. A web-site containing registration and user pages.
5. An API Gateway for securely interfacing client components with back end
services.

### Interaction
- A user first registers with the system through a user registration web page.
- Once registered, a user accesses the system by first authenticating through
the system's authentication API. This returns an expirable access token with
which the user passes to the remaing APIs used to access the system.
- With an access token obtained, a user can list, add, update, and delete
individual contact records through respective APIs.

## Design
### Registration
User information stored in the registration database includes:

- ctuser: the user's Cloudtacts login identifier            (max length 20)
- ctpass: the user's Cloudtacts login password (sha-256)    (length 64)
- ctprof: the user's profile name (displayed on site)       (max length 20)
- ctppic: the user's profile image key in object storage    (max length 52)
- uemail: the user's e-mail address                         (max length 50)
- atoken: the user's temporary access token                 (length 20)
- llogin: the user's last login timestamp                   (YYYYMMDDhhmmss)
- uvalid: the user's registration validation timestamp      (YYYYMMDDhhmmss)

When a user registers with the site as a new user, they're prompted for a user
identifier, password, and e-mail address, along with a profile name and an
optional profile picture in GIF, JPEG, or PNG format to display on the user
access page.

A user's login identifier, profile name, and e-mail address are all validated
to be unique within the system. Once validated, Cloudtacts:

1. generates a temporary access token with a 15 minute expiration, and
2. saves the user's information and temporary access token in the database, and
3. sends a confirmation e-mail to the provided e-mail address with a
confirmation link back to the authentication endpoint, and
4. starts a 15 minute validation timer for the new user.

The confirmation link and validation timer both contain references to the
user's login identifier, profile name, and temporary access token. If the
confirmation link is clicked by the user before the validation timer expires
then:

1. the validation timer is stopped and discarded, and
2. the user's information in the database is updated with a validation
timestamp to indicate successful validation by the user, and
3. the user's browser is redirected to the user access page.

However, if validation timer expires before the confirmation link is clicked
then the user's information is deleted from the database and the user directed
to return to the registration page to try again.

#### User Password
User passwords are stored in the database as sha-256 hashed values to prevent
discovery by third-parties.

#### User Profile Image
Users' optional profile images are saved to object storage with key pattern:
{ctuser}/{ctprof}/image.{ext}, where {ext} is the uploaded file name extension
(gif, jpeg/jpg, png). Object storage keys are saved with the user's information
in the database.

### Authentication
When a user authenticates with the application by signing in through the user
access page, they're prompted for their login identifier and password.

TODO

## Implementation
All back-end components are written in Go(lang). A registration and simple
contacts management page is written in PHP/HTML and JavaScript using the React
library. (_TBD_)

- **Registration and Authentication**

    A registration and authentication interface is developed as a pair of Cloud
("lambda") Functions written in Go(lang). These functions interface with a
user database in MySQL maintained on Cloud SQL.

- **User Access**

    A user access interface for listing, adding, updating, and deleting contact
records is developed as a Cloud Run container written in Go(lang), packaged
with Docker, and orchestrated with Google Kubernetes Engine (GKE). The
container interfaces with a Firestore document database where users' contact
records are kept.

- **User Interface**

    A web-based user interface is developed as an AppEngine stack containing
pages for user registration and database access, written in PHP/HTML and
JavaScript. The pages interface with the registration and authentication APIs
serviced by Cloud Functions and the access APIs serviced by Cloud Run,
utilizing API Gateway.

## Testing
TODO

## Deployment
TODO
