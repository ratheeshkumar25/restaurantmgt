**Restaurant Management**
**User Authentication: This section likely deals with creating user accounts, logging in, and managing user permissions.**

**Menu Management: This section most likely allows users to create menus, add items to menus, and update pricing.**

**Staff Management: This section likely allows users to add staff information, assign roles, and manage schedules.**

**Table Management: This section likely allows users to create a digital floor plan of the restaurant, assign tables to reservations, and track the status of tables (e.g., occupied, dirty, clean).**

**Order and Invoice Management: This section likely allows users to take orders, generate bills, and track order history.**

**Reviews: This section likely allows users to view and manage customer reviews.**

﻿

**User**
This folder encompasses functionalities related to user authentication, including user login and login verification processes. It is responsible for managing user login sessions, verifying user credentials, and ensuring secure access to the application.EndFragment

﻿

POST
login
https://rerarestaurantmanagement.shop/users/login
﻿

Authorization
Bearer Token
Token
•••••••
Body
raw (json)
json
{
    "phone":"+919353306805"
   
}
GET
home
https://rerarestaurantmanagement.shop/users
﻿

Body
raw (json)
json
{
    "phone":"+917356910125"
}
POST
login/Verify
https://rerarestaurantmanagement.shop/users/login/verify
﻿

Body
raw (json)
json
{
    "phone": "+919353306805",
    "otp":"837694"
}
**Admin**
The API folder hosts functionalities dedicated to admin authentication, token generation, and logout processes. It is designed to ensure seamless interaction with the server, guaranteeing proper responses for each operation as needed.EndFragment

﻿

POST
login
https://rerarestaurantmanagement.shop/admin/login
﻿

Query Params
Body
raw (json)
json
{
    "Username" : "admin",
    "Password" : "admin12345"
}
POST
logout
https://rerarestaurantmanagement.shop/admin/logout
﻿

**MenuMangement**
The API folder oversees essential tasks related to managing menu items, including fetching existing items, adding new items, updating item details, and removing items from the menu. Its primary focus is on facilitating seamless communication with the server and delivering clear and concise responses for each action performedEndFragment

﻿

User
﻿

GET
**getMenu**
http://rerarestaurantmanagement.shop/users/menu/1
﻿

Authorization
Bearer Token
Token
•••••••
GET
**menuListUsers**
http://rerarestaurantmanagement.shop/users/menulist
﻿

Authorization
Bearer Token
Token
•••••••
Admin
﻿

POST
**addItems**
https://rerarestaurantmanagement.shop/admin/menu/add
﻿

Authorization
Bearer Token
Token
•••••••
Body
raw (json)
View More
json
{
    "Food_ID":     8,
    "category":   "Bevarage",
    "name":       "Milk-Shake",
    "price":       110.55,
    "foodimage": "https://www.yummytummyaarthi.com/wp-content/uploads/2022/07/chicken-clear-soup-1.jpg",
    "duration": "25 min",
    "TableID" : 6
}
GET
**menuListAdmin**
https://rerarestaurantmanagement.shop/admin/menuList
﻿

Authorization
Bearer Token
Token
•••••••
PUT
**editItems**
https://rerarestaurantmanagement.shop/admin/menu/4
﻿

Authorization
Bearer Token
Token
•••••••
Body
raw (json)
View More
json

    {
    "ID"    :4,
    "category":   "Chicken-Soup",
    "name":       "Soup",
    "price":       30.55,
    "foodimage": "https://www.yummytummyaarthi.com/wp-content/uploads/2022/07/chicken-clear-soup-1.jpg",
    "duration":   "2024-03-18T12:00:00Z",
    "TableID" : 6
    }
DELETE
**deleteItems**
https://rerarestaurantmanagement.shop/admin/menu/4
﻿

Authorization
Bearer Token
Token
•••••••
**StaffManagement**
The API folder manages tasks like fetching,adding, updating, assigning tables, and removing staff. It ensures smooth communication with the server and provides clear responses for each action.

﻿

GET
**getStaffById**
https://rerarestaurantmanagement.shop/admin/staff/1
﻿

Authorization
Bearer Token
Token
•••••••
PUT
**updateStaff**
https://rerarestaurantmanagement.shop/admin/staff/1
﻿

Authorization
Bearer Token
Token
•••••••
Body
raw (json)
json
{
  
    "staffname":"Govind",
    "staffrole":"Waiter",
    "salary":19000
}
DELETE
**removeStaff**
https://rerarestaurantmanagement.shop/admin/staff/2
﻿

Authorization
Bearer Token
Token
•••••••
POST
**addStaff**
https://rerarestaurantmanagement.shop/admin/staff/add
﻿

Authorization
Bearer Token
Token
•••••••
Body
raw (json)
json
{
    "staffname":"Raghav",
    "staffrole":"Waiter",
    "salary":5500

}
GET
**getStaff**
https://rerarestaurantmanagement.shop/admin/staff
﻿

Authorization
Bearer Token
Token
•••••••
**TableManagement**
The API folder orchestrates a variety of tasks including fetching data, reserving tables, updating reservations, and canceling reservations. Additionally, it enables administrators to add, update, and remove tables along with their respective capacities. Its core objective is to facilitate seamless communication with the server while furnishing transparent responses for every action undertakenEndFragment

﻿

**Userside**
﻿

GET
**cancelReservation**
https://rerarestaurantmanagement.shop/users/cancelreservation/4
﻿

Authorization
Bearer Token
Token
•••••••
GET
**userViewTables**
https://rerarestaurantmanagement.shop/users/table
﻿

Authorization
Bearer Token
Token
•••••••
POST
**reserveTable**
http://rerarestaurantmanagement.shop/users/reservation
﻿

Authorization
Bearer Token
Token
•••••••
Body
raw (json)
json
{
    "date": "2024-04-24T17:15:00Z",
    "numberOfGuest":5,
    "startTime": "2024-04-24T19:50:00Z",
    "endTime": "2024-04-24T20:10:00Z"
}


GET
**searchtable**
http://rerarestaurantmanagement.shop/users/searchreservation
﻿

Authorization
Bearer Token
Token
•••••••
Body
raw (json)
json
{
    "startTime": "2024-04-13T19:50:00Z",
    "endTime": "2024-04-13T20:10:00Z",
    "numberofGuest": 1
}
PUT
**moveReservation**
http://rerarestaurantmanagement.shop/users/movereservation/3
﻿

Authorization
Bearer Token
Token
•••••••
Body
raw (json)
json
{
    "dateTime": "2024-04-12T08:00:00Z",
    "numberOfGuest": 5,
   "startTime": "2024-04-17T19:50:00Z",
    "endTime": "2024-04-17T20:10:00Z"

}
Adminside
﻿

GET
**adminViewTables**
http://rerarestaurantmanagement.shop/admin/table
﻿

Authorization
Bearer Token
Token
•••••••
DELETE
clearTable
localhost:3000/admin/table/6
﻿

Authorization
Bearer Token
Token
•••••••
POST
**addTable**
http://rerarestaurantmanagement.shop/admin/table/add
﻿

Authorization
Bearer Token
Token
•••••••
Body
raw (json)
json
{
    "capacity":8,
    "availability":true
}
PUT
**updateTable**
https://rerarestaurantmanagement.shop/admin/table/1
﻿

Authorization
Bearer Token
Token
•••••••
Query Params
Body
raw (json)
json
{
    "capacity":4,
    "availability":true
}
**Order and InvoiceManagement**
The API folder oversees a range of tasks such as fetching data, placing orders, updating information, generating invoices, paying invoices, and canceling orders and invoices. Its primary function is to facilitate seamless communication with the server and furnish transparent responses for every action initiated. In addition to these functionalities, the administrative side has the capability to generate sales reports, revenue reports, and assess employee performance, along with tracking the total number of orders processed.EndFragment

﻿

**User**
﻿

POST
**placeOrder**
https://rerarestaurantmanagement.shop/users/placeorder/invoice/
﻿

Authorization
Bearer Token
Token
•••••••
Body
raw (json)
View More
json
{
  "items": [
    {
      "itemID":1,
      "quantity": 6
    },
    {
      "itemID": 3,
      "quantity": 4
    }
  ],
  "paymentMethod": "online",
  "email":"ratheeshgopinadhkumar@gmail.com"
}
PUT
**updateOrder**
https://rerarestaurantmanagement.shop/users/updateorder/18
﻿

Authorization
Bearer Token
Token
•••••••
Body
raw (json)
View More
json
{
  "items": [
    {
      "itemID": 5,
      "quantity": 2
    },
    {
      "itemID": 1,
      "quantity":2
    }
  ],
  "email":"ratheeshgopinadhkumar@gmail.com",
  "paymentMethod": "online"
}
     
POST
**payInvoice**
https://rerarestaurantmanagement.shop/users/payinvoice/1
﻿

Authorization
Bearer Token
Token
•••••••
GET
**cancelOrder**
https://rerarestaurantmanagement.shop/users/cancelorder/2
﻿

Authorization
Bearer Token
Token
•••••••
Admin
﻿

**SaleReport**
﻿

GET
**totalOrder**
http://rerarestaurantmanagement.shop/admin/totalorder
﻿

Authorization
Bearer Token
Token
•••••••
GET
**totalSales**
http://rerarestaurantmanagement.shop/admin/sales
﻿

Authorization
Bearer Token
Token
•••••••
GET
**employeePerformance**
http://rerarestaurantmanagement.shop/admin/employeeperformance
﻿

Authorization
Bearer Token
Token
•••••••
GET
**revenueReport**
http://rerarestaurantmanagement.shop/admin/revenue
﻿

Authorization
Bearer Token
Token
•••••••
GET
**adminInvoiceViews**
http://rerarestaurantmanagement.shop/admin/invoices/2/pdf/
﻿

Authorization
Bearer Token
Token
•••••••
Query Params
id
3
**Reviews**
The API folder is responsible for handling tasks related to managing user reviews and feedback, including fetching existing reviews and feedback from the server, as well as allowing users to post new reviews and feedback.EndFragment

﻿

GET
**ratingView**
http://rerarestaurantmanagement.shop/users/rating
﻿

Authorization
Bearer Token
Token
•••••••
POST
**feedBack**
http://localhost:3000/users/rating
﻿

Authorization
Bearer Token
Token
•••••••
Body
raw (json)
json
{
    "name":"Ratheesh",
    "suggestion":"Over all awesome service",
    "rating":9
}
