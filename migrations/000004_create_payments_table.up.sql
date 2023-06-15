CREATE TABLE payments (
                          id INT AUTO_INCREMENT PRIMARY KEY,
                          stripe_payment_id VARCHAR(255) NOT NULL,
                          user_id VARCHAR(15) NOT NULL,
                          amount DECIMAL(10, 2) NOT NULL,
                          currency VARCHAR(3) NOT NULL,
                          payment_method VARCHAR(255) NOT NULL,
                          parent_id VARCHAR(15) NOT NULL,
                          status INT NOT NULL DEFAULT 0,
                          created_at TIMESTAMP NOT NULL,
                          updated_at TIMESTAMP NOT NULL
);