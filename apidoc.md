# Api Documentation Loan Process

***POST /loan-applications***
----
Create new loan application
* **Headers**  
  *Content-Type:* `multipart/form-data`
* **Data Params**
  ```
    {
      full_name: string
      ktp_number: string
      gender: string
      date_of_birth: date
      address: string
      phone_number: string
      email: string
      nationality: string
      address_province: string
      ktp_image: image
      selfie_image: image
      loan_amount: number
      tenor: number
    }
  ```
* **Success Response:**
* **Code:** 201
  * **Content:**
    ```json
    {
        "status": "CREATED",
        "code": 201,
        "message": "Success create loan application"
    }
    ```
* **Error Response:**
    * **Code:** 422
        * **Content:**
        ```json
        {
            "status": "UNPROCESSABLE_ENTITY",
            "code": 422,
            "errors": [
                {
                    "ktp_number": "The ktp number already exist."
                }
            ]
        }
        ```
  * **Code:** 400
      * **Content:**
      ```json
      {
          "status": "BAD_REQUEST",
          "code": 400,
          "message": "The loan application daily limit exceeded"
      }
      ```

***Get /loan-applications***
----
Get list loan application
* **Headers**  
  *Content-Type:* `application/json`
* **Success Response:**
* **Code:** 200
  * **Content:**
    ```json
    {
        "status": "OK",
        "code": 200,
        "data": [
        {
          "customer_id": 1,
          "full_name": "Farhan",
          "ktp_number": "1234567890123123",
          "email": "farhan@gmail.com",
          "loan_amount": 1000000,
          "tenor": 3,
          "status": "accepted"
        }
      ]
    }
    ```

***Get /loan-applications/:customer_id***
----
Get loan application by id
* **URL Params**  
  *Required:* `customer_id=[uint]`
* **Headers**  
  *Content-Type:* `application/json`
* **Success Response:**
  * **Code:** 200
    * **Content:**
      ```json
      {
          "status": "OK",
          "code": 200,
          "data": {
            "full_name": "Farhan",
            "ktp_number": "1234567890123123",
            "gender": "Male",
            "date_of_birth": "2001-01-01",
            "address": "Jl. Test",
            "phone_number": "081234567890",
            "email": "farhan@gmail.com",
            "nationality": "INDONESIA",
            "address_province": "SUMATERA UTARA",
            "ktp_image": "https://loan.com/assets/ktp",
            "selfie_image": "https://loan.com/assets/selfie",
            "loan_amount": 1000000,
            "tenor": 3,
            "status": "accepted"
          }
      }
      ```
* **Error Response:**
    * **Code:** 404
      * **Content:**
        ```json
        {
            "status": "NOT_FOUND",
            "code": 404,
            "message": "The loan application not found"
        }
        ```
***POST /loan-applications/:customer_id/reapply***
----
Reapply new loan application
* **URL Params**  
  *Required:* `customer_id=[uint]`
* **Headers**  
  *Content-Type:* `application/json`
* **Success Response:**
* **Code:** 200
  * **Content:**
    ```json
    {
        "status": "OK",
        "code": 200,
        "message": "Success reapply loan application"
    }
    ```
* **Error Response:**
  * **Code:** 404
    * **Scenario:** Loan application limit
    * **Content:**
      ```json
      {
          "status": "BAD_REQUEST",
          "code": 400,
          "message": "The loan application daily limit exceeded"
      }
      ```
    * **Scenario:** Customer already have accepted loan
    * **Content:**
      ```json
      {
          "status": "BAD_REQUEST",
          "code": 400,
          "message": "The customer already have accepted loan"
      }
      ```