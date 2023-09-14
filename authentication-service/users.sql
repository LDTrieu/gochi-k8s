CREATE TABLE public.users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email character varying(255),
    first_name character varying(255),
    last_name character varying(255),
    password character varying(60),
    user_active integer DEFAULT 0,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

