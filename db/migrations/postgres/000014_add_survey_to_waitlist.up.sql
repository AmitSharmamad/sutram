-- Create survey table with enum types
CREATE TYPE nature_of_work AS ENUM ('Individual', 'Company', 'Proprietary_Trading_Firm', 'Asset_Management', 'Family_Office', 'Other');
CREATE TYPE trading_type AS ENUM ('Equity', 'FNO', 'Both');
CREATE TYPE experience_level AS ENUM ('Beginner', 'Intermediate', 'Pro');
CREATE TYPE work_type AS ENUM ('Personal', 'Firm');

CREATE TABLE survey (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    nature_of_work nature_of_work NOT NULL,
    trading_type trading_type NOT NULL,
    experience_level experience_level NOT NULL,
    work_type work_type NOT NULL,
    age INTEGER NOT NULL CHECK (age > 0 AND age <= 120),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_survey_email ON survey(email);
CREATE INDEX idx_survey_created_at ON survey(created_at);
CREATE INDEX idx_survey_nature_of_work ON survey(nature_of_work);
CREATE INDEX idx_survey_trading_type ON survey(trading_type);
CREATE INDEX idx_survey_experience_level ON survey(experience_level);

-- Update waitlist table to include survey reference
ALTER TABLE waitlist ADD COLUMN survey_id UUID REFERENCES survey(id) ON DELETE SET NULL;
CREATE INDEX idx_waitlist_survey_id ON waitlist(survey_id);
