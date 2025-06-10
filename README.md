# :pray: Prayer Buddies - Backend :pray: ![CI Badge](https://github.com/megarage9000/Prayer-Buddies/actions/workflows/ci.yaml/badge.svg)![Development Deployment Badge](https://github.com/megarage9000/Prayer-Buddies/actions/workflows/deploy-dev.yaml/badge.svg)![Production Deployment Badge](https://github.com/megarage9000/Prayer-Buddies/actions/workflows/deploy-prod.yaml/badge.svg)

Prayer Buddies - Backend is the backend system used for the Prayer Buddies application. This project handles user requests and server interactions with the Prayer Buddies database. Features include:

- Authentication using **JWT tokens** and password hashing using bcrypt
- ORM function using generated **SQLC** to communicate with the Prayer Buddies **PostgreSQL** Database
- **CRUD REST endpoints** that handle:
    - Registering and logging users in, responding with a JWT
    - Sending prayer requests, uploaded to Prayer Buddies PostgreSQL Database
    - Listing sent / received prayer requests from Prayer Buddies PostgreSQL Database

- Deployment Pipelines to deploy to [Google Cloud Platform](https://cloud.google.com/gcp?utm_source=bing&utm_medium=cpc&utm_campaign=na-CA-all-en-dr-bkws-all-all-trial-e-dr-1710134&utm_content=text-ad-none-any-DEV_c-CRE_-ADGP_Desk+%7C+BKWS+-+EXA+%7C+Txt-Core-General+GCP-KWID_43700063341856130-kwd-76553702185444:loc-32&utm_term=KW_gcp-ST_gcp&gclid=b572f1c855dc1cbe0aa1998f82937b07&gclsrc=3p.ds&msclkid=b572f1c855dc1cbe0aa1998f82937b07&hl=en) in both development and production environments
- Interactions with PostgreSQL hosted on [Neon PostgreSQL](https://neon.com/)
# Requirements

- Linux / UNIX operating system (This was run in WSL2)
- Go version 1.23.x
- Postgre SQL (This was using PostgreSQL 16.9)

# How To Run

To be written!

# :hammer: TODO Checklist :hammer:

- [ ] Documentation on how the backend works
- [ ] Implement refresh tokens
- [ ] Model Friends list per user
- [ ] Implement system to add friends / remove friends
- [ ] Implement search for users
- [ ] (Optional) Notifications (i.e. email / sms)
- [ ] (Optional) Send "prayed" notifications on prayed cards

