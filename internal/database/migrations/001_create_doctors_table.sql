--
CREATE TABLE IF NOT EXISTS doctors (
    id UUID DEFAULT uuid_generate_v7() PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    specialty_id UUID NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    avatar_url TEXT,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_doctors_specialty_id FOREIGN KEY (specialty_id) REFERENCES specialties(id) ON DELETE RESTRICT ON UPDATE CASCADE
);

--
CREATE INDEX IF NOT EXISTS idx_doctors_specialty_id ON doctors(specialty_id);

--
CREATE TRIGGER update_doctors_updated_at
    BEFORE UPDATE ON doctors
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
