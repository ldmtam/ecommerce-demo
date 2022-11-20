## cURL
### Create new product
```bash
curl --location --request POST 'localhost:3000/api/v1/products' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "name": "Ultraboost 22 shoes",
        "price": 250
    }'
```

```bash
curl --location --request POST 'localhost:3000/api/v1/products' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "name": "Ultraboost 4DFWD shoes",
        "price": 300
    }'
```

```bash
curl --location --request POST 'localhost:3000/api/v1/products' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "name": "Stan Smith shoes",
        "price": 200
    }'
```

### Get product by id
```bash
curl --location --request GET 'localhost:3000/api/v1/products/1' \
    --header 'Content-Type: application/json' \
    --header 'Cookie: user_id=123' \
    --data-raw ''
```

```bash
curl --location --request GET 'localhost:3000/api/v1/products/2' \
    --header 'Content-Type: application/json' \
    --header 'Cookie: user_id=456' \
    --data-raw ''
```

### Seach product by names
```bash
curl --location --request GET 'localhost:3000/api/v1/products/seachByName/boost' \
    --header 'Cookie: user_id=123' \
    --data-raw ''
```

```bash
curl --location --request GET 'localhost:3000/api/v1/products/seachByName/shoe' \
    --header 'Cookie: user_id=123' \
    --data-raw ''
```

### Get customer activities
```bash
curl --location --request GET 'localhost:3000/api/v1/customer_activities/123' \
    --data-raw ''
```

```bash
curl --location --request GET 'localhost:3000/api/v1/customer_activities/456' \
    --data-raw ''
```

```bash
curl --location --request GET 'localhost:3000/api/v1/customer_activities/123/actions/VIEW_PRODUCT' \
    --data-raw ''
```