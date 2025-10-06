🧩 API Endpoints (MVP Phase 1)
Method	Endpoint	Description
POST	/api/v1/register	Create user
POST	/api/v1/login	Authenticate user
GET	/api/v1/listings	Fetch listings (filters: city, price, type)
POST	/api/v1/listings	Create new listing (agent only)
GET	/api/v1/listings/:id	Get listing details
POST	/api/v1/alerts	Create alert
GET	/api/v1/alerts	Get user alerts
POST	/api/v1/favorites	Save listing as favorite
GET	/api/v1/favorites	Get all favorites