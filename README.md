# Go Example Project

This projects aim to practice following topics,

- Message Broker (w/ RabbitMQ)
- Caching (w/ Redis)
- API (on both server & client side)
- Go 
 - Go routine
 - HTTP Framework (w/ Echo)
- Database
 - MongoDB
 - Redis (cache)
- Microservice patterns 
 - Event sourcing


# A Order Service

คิดซะว่าเป็นส่วนหนึ่งของระบบ E-Commerce อะไรสักอย่าง ที่เป็นระบบใหญ่ เลยจำเป็นต้องแยกออกมาเป็น Service เดี่ยวๆ ที่ดูเเลเกี่ยวกับการจัดการ Order เท่านั้น

ขั้นตอนการทำงานเป็นแบบนี้ 

https://lucid.app/lucidchart/47d73a0e-a713-4136-be28-524dbbf4e89a/edit?viewport_loc=-1691%2C-1560%2C1984%2C1131%2C0_0&invitationId=inv_b84fbe3a-6396-4327-aea6-c820a7361b46

## Order
```
order_id    string
status      string
created_at  string
updated_at  string
```

## Order Status
```
order_id    string
status      string
created_at  string
remark      string
```

เก็บอยู่ในฐานข้อมูล 2 ก้อน (2 collection) คือ

- order เก็บสถานะปัจจุบันของ order
- order_status เก็บสถานะทุกอันของ order นั้นๆ เพราะ order จะสามารถเข้ามาซ้ำๆได้ ถ้าเข้ามาซ้ำก็จะสร้างใหม่เลย ไม่ update ของเก่า เเต่จะอ้างอิงจากของใหม่ที่สุดเท่านั้น


Status
```
ประกอบด้วย
- PENDING
- PAID (จ่ายแล้ว รอส่งสินค้า)
- SHIPPED (ส่งแล้ว ลูกค้ายังไม่ยืนยัน)
- COMPLETE
- CANCELED
```

- สาเหตุที่ให้แยกเป็น `order` กับ `order_status` เพราะเก็บ order_status เป็นประวัติด้วยว่ามีการอัพเดทเมื่อไร เเละเป็นการใช้ pattern sourcing เผื่อมีการดึงไปใช้งานในอนาคตด้วย
- สาเหตุที่ให้แยกเป็น API ในการอัพเดท เนื่องจากจะให้ลองใช้ Pattern CQRS

# Reference

https://medium.com/@chatthanajanethanakarn/cqrs-%E0%B8%89%E0%B8%9A%E0%B8%B1%E0%B8%9A%E0%B8%AD%E0%B9%88%E0%B8%B2%E0%B8%99%E0%B8%9A%E0%B8%99%E0%B8%A3%E0%B8%96%E0%B9%84%E0%B8%9F%E0%B8%9F%E0%B9%89%E0%B8%B2-dbf44e9a2dc1

https://medium.com/@chatthanajanethanakarn/event-sourcing-%E0%B9%81%E0%B8%9A%E0%B8%9A%E0%B8%AA%E0%B8%B1%E0%B9%89%E0%B8%99%E0%B9%86-%E0%B9%80%E0%B8%99%E0%B8%B7%E0%B9%89%E0%B8%AD%E0%B9%86-ea4fb24158c6

