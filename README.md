# Dev

    API for a Developer Connect Portal.
    A Portal to get the Developers Connected, Post their thoughts and get networked easily.
### Front End:

    The Project has its templates and all the static source in the templates folder. 
    The API's has not been Integrated in the Frontend yet. 

### Back End

#### Version of API :  V1.0
    
# API Documentation
 All responses come in standard JSON. All requests must include a `content-type` of `application/json` and the body must be valid JSON.


### Response Codes
```
200: Success
201: Created
400: Bad request
401: Unauthorized
404: Cannot be found
405: Method not allowed
50x: Server Error
```
### Error and Success Message Example

```json
  {
    "error":"message" 
  }
  
  {
      "success":"message",
      "any-data-type": "data-sent-in response"  //w.r.t the API
  }
```

## SignUp

**You send:**  You send the details required to signup.

**You get:** An `Error-Message` or a `Success-Message` depending on the status of the account created on DevConnect

**Endpoint:** 
     /Dev/signup

**Authorization Token:** Not required

**Request:**
`POST HTTP/1.1`

```json
Accept: application/json
Content-Type: application/json
Content-Length: xy

{   
    "name": "abc",
    "email": "foo",
    "password": "1234567",
    "cpassword": "1234567"
}
```

**Successful Response:**
```json
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: xy

{
   "success":"created",
}
```

## Login

**You send:**  Your  login credentials.

**You get:** An `API-Token` and a `Success-Message` with which you can make further actions.

**Endpoint:** 
     /Dev/login

**Authorization Token:** Not required

**Request:**
`POST HTTP/1.1`
```json
Accept: application/json
Content-Type: application/json
Content-Length: xy

{
    "email": "foo",
    "password": "1234567" 
}
```

**Successful Response:**
```json
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: xy

{
   "success":true,
   "token": "e3b0...................",
}
```
## Get Profile

**You send:**  unique user ID.

**You get:** A `Profile data` or a `Error-Message` with which you can make further actions.

**Endpoint:** 
     /Dev/Profile/{id}

**Authorization Token:** Bearer token required

**Request Param**
   `GET HTTP/1.1` 
   ```
   id  : Unique ID of the user 
   ```


**Successful Response:**
```json
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: xy

{
   "uuid":"",
   "email":"",
   "name":"",
   "status":"",
   "org":"",
   "website":"",
   "location":"",
   "skills":[],
   "gitname":"",
   "bio":"",
   "social":{},
   "edu":[
            {
            "school":"",
            "degree":"",
            "fields":"",
            "from":"",
            "to":"",
            "achievements":""
            }
        ],
   "exp":[
            {
            "title":"",
            "org":"",
            "location":"",
            "from":"",
            "to":"",
            "description":""
            }
       ]
}
```

## Dashboard

**You send:**  nothing

**You get:** A `Profile data  i.e the Education and Experience` or a `Error-Message` with which you can make further actions.

**Endpoint:** 
     /Dev/Dashboard

**Authorization Token:** Bearer token required

**Request**
   `GET HTTP/1.1` 
    
    Nil

**Successful Response:**
```json
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: xy

{
    "name":"",
    "education":[
            {
               "school":"",
               "degree":"",
               "fields":"",
               "from":"",
               "to":"",
               "achievements":""
            }
          ],
    "experience":[
             {
               "title":"",
               "org":"",
               "location":"",
               "from":"",
               "to":"",
               "description":""
            }
          ]   
}
```

## Update Profile

**You send:**  The data to be updated in the profile.

**You get:** A `Success-Message` or a `Error-Message` with which you can make further actions.

**Endpoint:** 
     /Dev/profile/update

**Authorization Token:** Bearer token required

**Request**
   `POST HTTP/1.1` 
```json
Accept: application/json
Content-Type: application/json
Content-Length: xy

{
    "status":"",
    "org":"",
    "website":"",
    "location":"",
    "skills":[],
    "gitname":"",
    "bio":"",
    "social":{
       "twitter":"",
       "facebook":"",
       "linkedin":"",
       "instagram":"",
       "youtube":""
             }
}
```

**Successful Response:**
```json
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: xy

{
   "success":"updated"
}
```

## Update Education

**You send:**  The data to be added or updated in the education.

**You get:** A `Success-Message` or a `Error-Message` with which you can make further actions.

**Endpoint:** 
     /Dev/profile/add/Experience

**Authorization Token:** Bearer token required

**Request**
   `POST HTTP/1.1` 
```json
Accept: application/json
Content-Type: application/json
Content-Length: xy
 {
    "school":"",
    "degree":"",
    "fields":"",
    "from":"",
    "to":"",
    "achievements":""
 }
```

**Successful Response:**
```json
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: xy

{
   "success":"updated"
}
```

## Update Experience

**You send:**  The data to be added or updated in the experience.

**You get:** A `Success-Message` or a `Error-Message` with which you can make further actions.

**Endpoint:** 
     /Dev/profile/add/Experience

**Authorization Token:** Bearer token required

**Request**
   `POST HTTP/1.1` 
```json
Accept: application/json
Content-Type: application/json
Content-Length: xy
 {
    "title":"",
    "org":"",
    "location":"",
    "from":"",
    "to":"",
    "description":""
 }
``` 

**Successful Response:**
```json
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: xy

{
   "success":"updated"
}
```

## Developers Data

**You send:**  Nothing

**You get:** The Data of all the developers.

**Endpoint:** 
     /Dev/Developers

**Authorization Token:** Bearer token required

**Request**
   `GET HTTP/1.1` 
    
    nil

**Successful Response:**
```json
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: xy

[
  {
    "uid":"",
    "name":"",
    "email":"",
    "status":"",
    "org":"",
    "website":"",
    "location":"",
    "skills":[],
    "gitname":"",
    "bio":"",
    "social":{
       "twitter":"",
       "facebook":"",
       "linkedin":"",
       "instagram":"",
       "youtube":""
             },
    "edu":[
            {
               "school":"",
               "degree":"",
               "fields":"",
               "from":"",
               "to":"",
               "achievements":""
            }
          ],
    "exp":[
             {
               "title":"",
               "org":"",
               "location":"",
               "from":"",
               "to":"",
               "description":""
            }
          ]      
  },
  {
     .....
  }
]
```

## Create Post

**You send:**  The data to be posted as Post.

**You get:**   A `Success-Message` or a `Error-Message` with which you can make further actions also the Total Posts .

**Endpoint:** 
     /Dev/Post

**Authorization Token:** Bearer token required

**Request**
   `GET , POST HTTP/1.1`  
```json
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: xy

   {
      "text":""
   }
``` 

**Successful Response:**
```json
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: xy

for POST Request:

{
   "success":"created"
}

for GET Request:

[
   {
      "id":"",
      "username":"",
      "email":"",
      "text":"",
      "comment":[],
      "likes": 0,  
      "dislikes": 0
   }
   {
      ....
   }
]
```

## Comment on Post

**You send:**  The `id` of the post and the comment.

**You get:**   A `All comments on the post` or a `Error-Message` with which you can make further actions  .


**Endpoint:** 
     /Dev/Post/comment/{id}

**Authorization Token:** Bearer token required

**Request Param**
   `GET HTTP/1.1`  
    
    id : ID of the Post

**Request Body**

```json
Content-Type: application/json
Content-Length: xy

{
   "text":"comment"
}
```
  
**Successful Response:**
```json
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: xy

//list of comments
[
   "hgfdcvbn",
   "huytfvbjk",
   ...

]

```

## Like Post

**You send:**  Post ID.

**You get:** A `Success-Message` or a `Error-Message` with which you can make further actions.

**Endpoint:** 
     /Dev/like/{id}

**Authorization Token:** Bearer token required

**Request Param**
   `GET HTTP/1.1` 
   ```
   id : ID of the Post
   ```

**Successful Response:**
```json
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: xy

{
   "success":"updated"
}
```

## Dislike Post

**You send:**  Post ID.

**You get:** A `Success-Message` or a `Error-Message` with which you can make further actions.

**Endpoint:** 
     /Dev/dislike/{id}

**Authorization Token:** Bearer token required

**Request Param**
   `GET HTTP/1.1` 
   ```
   id : ID of the Post 
   ```

**Successful Response:**
```json
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: xy

{
   "success": "updated"
}
```

**Failed Response for All API's:**

```json 
HTTP/1.1 400 Bad Request
Content-Type: application/json
Content-Length: xy

{
   "error":"message"
}
``` 

## Environment variables

  
  `MongoUrl` : to store the url of mongodb cluster or localhost.

  `BLOCKKEY` : to store encryption secret.

  `PORT`     : port on localhost to run application

