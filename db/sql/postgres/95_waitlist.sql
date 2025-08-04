-- Create waitlist table
CREATE TABLE waitlist (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    survey_id UUID REFERENCES survey(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on email for faster lookups
CREATE INDEX idx_waitlist_email ON waitlist(email);

-- Create index on created_at for analytics
CREATE INDEX idx_waitlist_created_at ON waitlist(created_at);

-- Create index on survey_id
CREATE INDEX idx_waitlist_survey_id ON waitlist(survey_id);
