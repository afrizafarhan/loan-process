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
            "error": "ktp image extension not allowed"
        }
        ```
  * **Code:** 400
      * **Content:**
      ```json
      {
          "status": "BAD_REQUEST",
          "code": 400,
          "error": "The loan application daily limit exceeded"
      }
      ```
      * **Content:**
      ```json
      {
          "status": "BAD_REQUEST",
          "code": 400,
          "error": "ktp number already exists"
      }
      ```
      * **Content:**
      ```json
      {
          "status": "BAD_REQUEST",
          "code": 400,
          "error": "email already exists"
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
          "id": 1,
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
          "error": "The loan application daily limit exceeded"
      }
      ```
    * **Scenario:** Customer already have accepted loan
    * **Content:**
      ```json
      {
          "status": "BAD_REQUEST",
          "code": 400,
          "error": "the customer already have accepted loan"
      }
      ```

***Get /v1/customers***
----
Get list customers
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
          "id": 1,
          "full_name": "Farhan",
          "ktp_number": "1234567890123123",
          "email": "farhan@gmail.com",
        }
      ]
    }
    ```


***Get /v1/customers/:id/detail**
Get detail data customer
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
          "id": 1,
          "full_name": "Farhan",
          "ktp_number": "1234567890123123",
          "date_of_birth": "2001-01-01",
          "address": "Jl. Test",
          "phone_number": "081234567890",
          "email": "farhan@gmail.com",
          "nationality": "indonesia",
          "address_province": "SUMATERA UTARA",
          "ktp_image": "http://loan.com/resources/ktp/image.png",
          "selfie_image": "http://loan.com/resources/selfie/image.png",
          "created_at": "2023-01-01",
          "updated_at": "2023-01-01"
        }
    }
    ```
* **Error Response:**
  * **Code:** 404
    * **Scenario:** Customer not found
    * **Content:**
      ```json
      {
          "status": "NOT_FOUND",
          "code": 404,
          "error": "customer not found"
      }
      ```

***Get /v1/customers/:id/loan-applications**
Get detail data customer
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
            "amount": 1000000,
            "tenor": 3,
            "status": "accepted",
            "created_at": "2023-01-01",
            "updated_at": "2023-01-01"
          },
      ]
    }
    ```
* **Error Response:**
  * **Code:** 404
    * **Scenario:** Customer not found
    * **Content:**
      ```json
      {
          "status": "NOT_FOUND",
          "code": 404,
          "error": "customer not found"
      }
      ```

***Get /v1/payment-installments/:loan_request_id***
----
Get list installment for accepted loan application
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
            "amount": 333333,
            "due_date": "2023-01-01",
            "status": "not_paid"
          },
          {
            "amount": 333333,
            "due_date": "2023-02-01",
            "status": "not_paid"
          },
          {
            "amount": 333333,
            "due_date": "2023-03-01",
            "status": "not_paid"
          }
      ]
    }
    ```
* **Error Response:**
  * **Code:** 404
    * **Scenario:** Accepted loan request not found
    * **Content:**
      ```json
      {
          "status": "NOT_FOUND",
          "code": 404,
          "error": "accepted loan request not found"
      }
      ```