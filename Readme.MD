# Forum Features

communication between users:
    posting
    liking
    commenting
    putting posts into categories
    post filters

## SQLite database
stores web content
    posts
    comments
    likes etc
must be controlled using queries and must use at least one of:
    select
    create
    insert queries

## Authentication
user profile creation
login feature, only logged in members can interact with posts, comments, likes and dislikes
can only have 1 opened session, use cookies
cookies must have expiration date, we can set that to whatever we like
use of UUID is a bonus

use of captcha to validate that user not a robot, will implement if time available

## User

e-mail
username
password

if e-mail or username taken, must return error, also checks for all credentials if they can log in or not
password to be encrypted when stored--- bonus

## User Comms
if registered {
    can post and comment
    when creating a post can use categories and filters to post in specified channels
    can like and dislike
    can view all of your like, dislikes, comments and posts
}

## filters
categories?

## Packages
this project can only use:

1. standard Go packages
2. sqlite3
3. bycrypt for encyption
4. UUID for cookies? 
5. Docker
6. Handle errors

## Errors

errors that have been identified and need to be accounted for:
    deleted posts and or comments, if trying to view them
    if empty posts or comment must return an error
    posts cannot be liked and dislkied, must accept newest option chosen
    when posting at least 1 category/filter must be selected
    if trying to log in to another session on same browser return an error as only 1 can be active at a time
    if !registered and trying to post, comment or like say they need to log in
    login and registration errors


### Forum type

no set type

### categories

### filters

### password requirements
will use industry standard password requirements but no special characters


### colour scheme

light & dark modes based on system theme

primary colour
secondary colour
tertiary colour

### layouts
the user will be greeted with a log in sctreen with options to log in, register or visit as guest

### User Profile
there is no specification here but at the minimum there will be 
    user name, and real name?
    user photo, selectable from stock photos for now
    a place to view all posts, comments, likes and dislkies. could they be viewed individually instead of all in one tab?
    like all likes, all dislkies etc

### special features 
here all the the ideas and or features we likeda bout different websites not all will be implemented

1. user to have ability to change colour scheme e.g. background colour and text colour or all 3 colours if 3 used?
2. look at tumblr--- where text is uploaded as an image?
3. voice messaging aka voice memos

 # Credits
    Karolis Milius
    Morayo
    Klaudia






ouststanding

user log out
forgot password
database
responsive web 




advanced features

notifications--- likes and dislikes, comments
edit/delete posts and comments

authenitction

login in with google and github

image upload
upload image into a post
sizes
format
where stored?

moderation
Guests

    These are unregistered-users that can neither post, comment, like or dislike a post. They only have the permission to see those posts, comments, likes or dislikes.

Users

    These are the users that will be able to create, comment, like or dislike posts.

Moderators

    Moderators, as explained above, are users that have a granted access to special functions :
        They should be able to monitor the content in the forum by deleting or reporting post to the admin
    To create a moderator the user should request an admin for that role

Administrators

    Users that manage the technical details required for running the forum. This user must be able to :
        Promote or demote a normal user to, or from a moderator user.
        Receive reports from moderators. If the admin receives a report from a moderator, he can respond to that report
        Delete posts and comments
        Manage the categories, by being able to create and delete them.


authentication

https
bcrypt 
uuid