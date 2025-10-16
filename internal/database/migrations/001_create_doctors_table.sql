-- Create doctors table
CREATE TABLE IF NOT EXISTS doctors (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    specialty_id int NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    avatar_url TEXT,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_doctors_specialty_id FOREIGN KEY (specialty_id) REFERENCES specialties(id) ON DELETE RESTRICT ON UPDATE CASCADE
);

-- Create index for faster queries
CREATE INDEX IF NOT EXISTS idx_doctors_specialty_id ON doctors(specialty_id);