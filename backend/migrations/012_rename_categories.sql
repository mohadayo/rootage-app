-- カテゴリ名の変更
UPDATE categories SET name = 'エンジニア基礎' WHERE name = '技術基礎';
UPDATE categories SET name = 'SES面談' WHERE name = '面談対策' OR name = 'SES面談対策';
UPDATE categories SET name = '現場マナー' WHERE name = 'ビジネスマナー';
UPDATE categories SET name = 'SES業界' WHERE name = 'SES理解' OR name = 'SES業界理解';
