# Dev

    API for a Developer Connect Portal.
    A Portal to get the Developers Connected, Post their thoughts and get networked easily.
### Front End:

    The Project has its templates and all the static source in the templates folder. 
    The API's have not been connected to the Frontend yet. 

### Back End

#### Version of API :  V1.0
    
# API Documentation Example
This API uses `POST` request to communicate and HTTP [response codes](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes) to indenticate status and errors. All responses come in standard JSON. All requests must include a `content-type` of `application/json` and the body must be valid JSON.

## Response Codes 
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

**Request Param**
   `GET HTTP/1.1` 
   ```
   id
   ```


**Successful Response:**
```json
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: xy

{
   
}
```

## Dashboard

**You send:**  nothing

**You get:** A `Profile data  i.e the Education and Experience` or a `Error-Message` with which you can make further actions.

**Endpoint:** 
     /Dev/Dashboard

**Request**
   `GET HTTP/1.1` 
    
    Nil

**Successful Response:**
```json
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: xy

{
   
}
```










**Failed Response for All API's:**
```json 
Content-Type: application/json
Content-Length: xy

{
   "error":"message"
}
``` 