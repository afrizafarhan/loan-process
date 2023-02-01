# Api Documentation Loan Process

***POST /v1/loan-applications***
----
Create new loan application
* **Headers**  
  *Content-Type:* `multipart/form-data`
* **Data Params**
  ```
    {
      full_name: string
      ktp_number: string
      gender: string (only accept male, female)
      date_of_birth: string (only accept format yy/mm/dd)
      address: string
      phone_number: string (length must bettween 10 and 13)
      email: string
      nationality: string (only accept value indonesia)
      address_province: string
      ktp_image: image (only accept mime type: jpg,png)
      selfie_image: image (only accept mime type: jpg,png)
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

***Get /v1/loan-applications***
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

***POST /v1/loan-applications/:customer_id/reapply***
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