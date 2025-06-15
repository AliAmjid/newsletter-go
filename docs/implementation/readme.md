# Implementation
This document describes the implementation of the backend part of the newsletter system. It provides information about the used packages, the database schema, and the individual tables that form the foundation of the application. The goal is to give an overview of the architecture and key components of the solution.

## Used packages
- [joho/godotenv](github.com/joho/godotenv)
- [go-chi/chi](https://github.com/go-chi/chi)

## Database

### Schema
![Database schema](../assets/database-schema.svg)

### Tables

#### user
| Name | Type | Settings | References |
| - | - | - | - |
| **id** | UUID | 🔑 PK, null | |
| **email** | TEXT | not null, unique | |
| **created_at** | TIMESTAMP | not null, default: now() | |
| **firebase_uid** | TEXT | null, unique | |

#### password_reset_tokens
| Name | Type | Settings | References |
| - | - | - | - |
| **token** | TEXT | 🔑 PK, null | |
| **user_id** | UUID | null | fk_password_reset_tokens_user_id_user | 
| **expires_at** | TIMESTAMP | not null | |
| **created_at** | TIMESTAMP | not null, default: now() | |

#### newsletter
| Name | Type | Settings | References |
| - | - | - | - |
| **id** | UUID | 🔑 PK, null | |
| **name** | TEXT | not null | |
| **description** | TEXT | null | |
| **owner_id** | UUID | not null | fk_newsletter_owner_id_user |
| **created_at** | TIMESTAMP | not null, default: now() | |

#### subscription
| Name | Type | Settings | References |
| - | - | - | - |
| **id** | UUID | 🔑 PK, null | |
| **newsletter_id** | UUID | not null | fk_subscription_newsletter_id_newsletter |
| **email** | TEXT | not null | |
| **token** | TEXT | not null, unique | |
| **confirmed_at** | TIMESTAMP | null | |
| **created_at** | TIMESTAMP | not null, default: now() | |

#### post
| Name | Type | Settings | References |
| - | - | - | - | 
| **id** | UUID | 🔑 PK, null | |
| **newsletter_id** | UUID | not null | fk_post_newsletter_id_newsletter |
| **title** | TEXT | not null | |
| **content** | TEXT | not null | |
| **published_at** | TIMESTAMP | null | |

#### post_delivery
| Name | Type | Settings | References |
| - | - | - | - |
| **id** | UUID | 🔑 PK, null | |
| **post_id** | UUID | not null | fk_post_delivery_post_id_post |
| **subscription_id** | UUID | not null | fk_post_delivery_subscription_id_subscription |
| **opened** | BOOLEAN | not null, default: false | |


### Relationships
- **newsletter to user**: many_to_one
- **password_reset_tokens to user**: many_to_one
- **post to newsletter**: many_to_one
- **subscription to newsletter**: many_to_one
- **post_delivery to post**: many_to_one
- **post_delivery to subscription**: many_to_one