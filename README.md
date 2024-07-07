# stage-two

HNG internship: stage two backend task

BACKEND Stage 2 Task: User Authentication & Organisation
Using your most comfortable backend framework of your choice, adhere to the following acceptance:
READ CAREFULLY!!!
Acceptance Criteria

Connect your application to a Postgres database server. (optional: you can choose to use any ORM of your choice if you want or not).
Create a User model using the properties below
NB: user id and email must be unique

```json
{
    "userId": "string" // must be unique
    "firstName": "string", // must not be null
    "lastName": "string" // must not be null
    "email": "string" // must be unique and must not be null
    "password": "string" // must not be null
    "phone": "string"
}
```

Provide validation for all fields. When there’s a validation error, return status code 422 with payload:
{
"errors": [
{
"field": "string",
"message": "string"
},
]
}

Using the schema above, implement user authentication
User Registration:

Implement an endpoint for user registration
Hash the user’s password before storing them in the database.

successful response: Return the payload with a 201 success status code.

User Login

Implement an endpoint for user Login.

Use the JWT token returned to access PROTECTED endpoints.
Organisation

A user can belong to one or more organisations
An organisation can contain one or more users.
On every registration, an organisation must be created.
The name property of the organisation takes the user’s firstName and appends “Organisation” to it. For example: user’s first name is John , organisation name becomes "John's Organisation" because firstName = "John" .
Logged in users can access organisations they belong to and organisations they created.
Create an organisation model with the properties below.
Organisation Model:

{
"orgId": "string", // Unique
"name": "string", // Required and cannot be null
"description": "string",
}

Endpoints:

[POST] /auth/register Registers a users and creates a default organisation Register request body:
{
"firstName": "string",
"lastName": "string",
"email": "string",
"password": "string",
"phone": "string",
}

Successful response: Return the payload below with a 201 success status code.

{
"status": "success",
"message": "Registration successful",
"data": {
"accessToken": "eyJh...",
"user": {
"userId": "string",
"firstName": "string",
"lastName": "string",
"email": "string",
"phone": "string",
}
}
}

Unsuccessful registration response:

{
"status": "Bad request",
"message": "Registration unsuccessful",
"statusCode": 400
}

[POST] /auth/login : logs in a user. When you log in, you can select an organisation to interact with
Login request body:

{
"email": "string",
"password": "string",
}

Successful response: Return the payload below with a 200 success status code.

{
"status": "success",
"message": "Login successful",
"data": {
"accessToken": "eyJh...",
"user": {
"userId": "string",
"firstName": "string",
"lastName": "string",
"email": "string",
"phone": "string",
}
}
}

Unsuccessful login response:

{
"status": "Bad request",
"message": "Authentication failed",
"statusCode": 401
}

[GET] /api/users/:id : a user gets their own record or user record in organisations they belong to or created [PROTECTED].
Successful response: Return the payload below with a 200 success status code.

{
"status": "success",
"message": "<message>",
"data": {
"userId": "string",
"firstName": "string",
"lastName": "string",
"email": "string",
"phone": "string"
}
}

[GET] /api/organisations : gets all your organisations the user belongs to or created. If a user is logged in properly, they can get all their organisations. They should not get another user’s organisation [PROTECTED].
Successful response: Return the payload below with a 200 success status code.

{
"status": "success",
"message": "<message>",
"data": {
"organisations": [
{
"orgId": "string",
"name": "string",
"description": "string",
}
]
}
}

[GET] /api/organisations/:orgId the logged in user gets a single organisation record [PROTECTED]
Successful response: Return the payload below with a 200 success status code.

{
"status": "success",
"message": "<message>",
"data": {
"orgId": "string", // Unique
"name": "string", // Required and cannot be null
"description": "string",
}
}

[POST] /api/organisations : a user can create their new organisation [PROTECTED].
Request body: request body must be validated

{
"name": "string", // Required and cannot be null
"description": "string",
}

Successful response: Return the payload below with a 201 success status code.

{
"status": "success",
"message": "Organisation created successfully",
"data": {
"orgId": "string",
"name": "string",
"description": "string"
}
}

Unsuccessful response:

{
"status": "Bad Request",
"message": "Client error",
"statusCode": 400
}

[POST] /api/organisations/:orgId/users : adds a user to a particular organisation
Request body:

{
"userId": "string"
}

Successful response: Return the payload below with a 200 success status code.

{
"status": "success",
"message": "User added to organisation successfully",
}

Unit Testing
Write appropriate unit tests to cover

Token generation - Ensure token expires at the correct time and correct user details is found in token.
Organisation - Ensure users can’t see data from organisations they don’t have access to.
End-to-End Test Requirements for the Register Endpoint
The goal is to ensure the POST /auth/register endpoint works correctly by performing end-to-end tests. The tests should cover successful user registration, validation errors, and database constraints.
Directory Structure:

The test file should be named auth.spec.ext (ext is the file extension of your chosen language) inside a folder named tests . For example tests/auth.spec.ts assuming I’m using Typescript
Test Scenarios:

It Should Register User Successfully with Default Organisation:Ensure a user is registered successfully when no organisation details are provided.
Verify the default organisation name is correctly generated (e.g., "John's Organisation" for a user with the first name "John").
Check that the response contains the expected user details and access token.
It Should Log the user in successfully:Ensure a user is logged in successfully when a valid credential is provided and fails otherwise.
Check that the response contains the expected user details and access token.
It Should Fail If Required Fields Are Missing:Test cases for each required field (firstName, lastName, email, password) missing.
Verify the response contains a status code of 422 and appropriate error messages.
It Should Fail if there’s Duplicate Email or UserID:Attempt to register two users with the same email.
Verify the response contains a status code of 422 and appropriate error messages.

### Step-by-Step Task Explanation

1. **Setup & Database Connection**

   - Choose a backend framework (e.g., Express.js for Node.js, Flask for Python, etc.).
   - Connect the application to a PostgreSQL database.
   - (Optional) Integrate an ORM (e.g., Sequelize for Node.js, SQLAlchemy for Python).
2. **User Model Creation**

   - Define a User model with the following properties:
     ```json
     {
       "userId": "string", // Unique
       "firstName": "string", // Required, not null
       "lastName": "string", // Required, not null
       "email": "string", // Unique, required, not null
       "password": "string", // Required, not null
       "phone": "string"
     }
     ```
   - Ensure `userId` and `email` are unique.
   - Add validation for all fields.
3. **Validation Handling**

   - Implement validation logic.
   - If validation fails, return a 422 status code with:
     ```json
     {
       "errors": [
         {
           "field": "string",
           "message": "string"
         }
       ]
     }
     ```
4. **User Authentication Implementation**

   - **User Registration**

     - Create a registration endpoint (`/auth/register`).
     - Hash the user’s password before storing it in the database.
     - On successful registration, return a 201 status code with:
       ```json
       {
         "status": "success",
         "message": "Registration successful",
         "data": {
           "accessToken": "eyJh...",
           "user": {
             "userId": "string",
             "firstName": "string",
             "lastName": "string",
             "email": "string",
             "phone": "string"
           }
         }
       }
       ```
     - On failure, return a 400 status code with:
       ```json
       {
         "status": "Bad request",
         "message": "Registration unsuccessful",
         "statusCode": 400
       }
       ```
   - **User Login**

     - Create a login endpoint (`/auth/login`).
     - Validate user credentials and generate a JWT token on successful login.
     - Return a 200 status code with:
       ```json
       {
         "status": "success",
         "message": "Login successful",
         "data": {
           "accessToken": "eyJh...",
           "user": {
             "userId": "string",
             "firstName": "string",
             "lastName": "string",
             "email": "string",
             "phone": "string"
           }
         }
       }
       ```
     - On failure, return a 401 status code with:
       ```json
       {
         "status": "Bad request",
         "message": "Authentication failed",
         "statusCode": 401
       }
       ```
5. **Organisation Management**

   - Define an Organisation model with the following properties:
     ```json
     {
       "orgId": "string", // Unique
       "name": "string", // Required, not null
       "description": "string"
     }
     ```
6. **Endpoints**

   - **[POST] /auth/register**

     - Registers a user and creates a default organisation with the user’s first name appended with "Organisation".
   - **[POST] /auth/login**

     - Logs in a user and returns a JWT token.
   - **[GET] /api/users/:id**

     - Retrieves a user’s record (PROTECTED endpoint).
     - Return 200 status code with user data.
   - **[GET] /api/organisations**

     - Retrieves all organisations the logged-in user belongs to or created (PROTECTED endpoint).
     - Return 200 status code with organisation data.
   - **[GET] /api/organisations/:orgId**

     - Retrieves a single organisation record (PROTECTED endpoint).
     - Return 200 status code with organisation data.
   - **[POST] /api/organisations**

     - Allows a user to create a new organisation (PROTECTED endpoint).
     - Return 201 status code with organisation data on success, 400 on failure.
   - **[POST] /api/organisations/:orgId/users**

     - Adds a user to a specific organisation (PROTECTED endpoint).
     - Return 200 status code on success.
7. **Unit Testing**

   - **Token Generation**

     - Ensure token expiration and correct user details in the token.
   - **Organisation Access**

     - Ensure users can’t see data from organisations they don’t have access to.
8. **End-to-End Testing for Register Endpoint**

   - **Test File Structure**

     - Name the test file `auth.spec.ext` inside a `tests` folder.
   - **Test Scenarios**

     - Successful user registration with default organisation.
     - Successful user login with valid credentials.
     - Registration failure for missing required fields.
     - Registration failure for duplicate email or userId.
