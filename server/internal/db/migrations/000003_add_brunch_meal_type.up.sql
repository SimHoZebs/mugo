DO $$ BEGIN
    ALTER TYPE meal_type ADD VALUE 'brunch' AFTER 'breakfast';
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;
