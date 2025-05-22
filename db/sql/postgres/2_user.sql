-- User Management Tables
CREATE TABLE IF NOT EXISTS "user" (
    id UUID DEFAULT gen_random_uuid(),
    name VARCHAR(255),
    username VARCHAR(255) UNIQUE,
    email VARCHAR(255) UNIQUE,
    phone_number VARCHAR(20) UNIQUE,
    password VARCHAR(255),
    verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    is_active BOOLEAN DEFAULT true,
    PRIMARY KEY (id)
);

-- Portfolio Management Tables
CREATE TABLE IF NOT EXISTS portfolio (
    id UUID DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE
);
