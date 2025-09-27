ALTER TABLE menu
ADD CONSTRAINT FK_MenuCategory
FOREIGN KEY (category_id) REFERENCES categories(id);