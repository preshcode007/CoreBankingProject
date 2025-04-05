import os
import psycopg2

def connect_db():
    try:
        db_host = os.environ.get("DATABASE_HOST")
        db_user = os.environ.get("DATABASE_USER")
        db_password = os.environ.get("DATABASE_PASSWORD")
        db_name = os.environ.get("DATABASE_NAME")
        db_port = os.environ.get("DATABASE_PORT")

        conn = psycopg2.connect(
            host=db_host,
            user=db_user,
            password=db_password,
            dbname=db_name,
            port=db_port
        )
        print("Successfully connected to PostgreSQL!")
        return conn

    except psycopg2.Error as e:
        print(f"Error connecting to database: {e}")
        return None