DROP TABLE IF EXISTS puzzle_progress;
DROP TABLE IF EXISTS puzzles;

CREATE TABLE puzzles (
    puzzle_id SERIAL PRIMARY KEY,
    category VARCHAR(255),
    url TEXT
);

ALTER TABLE puzzles ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Puzzles are viewable by all authenticated users"
ON puzzles FOR SELECT
TO authenticated
USING (true);

CREATE TABLE puzzle_progress (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES auth.users(id),
    puzzle_id INTEGER REFERENCES puzzles(puzzle_id),
    status VARCHAR(10) CHECK (status IN ('success', 'failed')),
    attempted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, puzzle_id)
);

-- RLS Policies
ALTER TABLE puzzle_progress ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can view own progress"
ON puzzle_progress FOR SELECT
USING (auth.uid() = user_id);

CREATE POLICY "Users can insert own progress"
ON puzzle_progress FOR INSERT
WITH CHECK (auth.uid() = user_id);

CREATE POLICY "Users can update own progress"
ON puzzle_progress FOR UPDATE
USING (auth.uid() = user_id);

-- All the 100 puzzles
INSERT INTO Puzzles (category, url) VALUES 
('Discovered Attacks', 'https://lichess.org/training/wVnOg'),
('Discovered Attacks', 'https://lichess.org/training/8P7Sd'),
('Discovered Attacks', 'https://lichess.org/training/1WXeO'),
('Discovered Attacks', 'https://lichess.org/training/BKGzE'),
('Discovered Attacks', 'https://lichess.org/training/SC6uI'),
('Discovered Attacks', 'https://lichess.org/training/kBDZ1'),
('Discovered Attacks', 'https://lichess.org/training/PnUWn'),
('Discovered Attacks', 'https://lichess.org/training/1yUnt'),
('Discovered Attacks', 'https://lichess.org/training/eVNd5'),
('Discovered Attacks', 'https://lichess.org/training/nobO7'),
('Discovered Attacks', 'https://lichess.org/training/StYI4'),
('Discovered Attacks', 'https://lichess.org/training/pDfoI'),
('Discovered Attacks', 'https://lichess.org/training/iivfc'),
('Discovered Attacks', 'https://lichess.org/training/ty8qp'),
('Discovered Attacks', 'https://lichess.org/training/vOZPN'),
('Discovered Attacks', 'https://lichess.org/training/sID53'),
('Discovered Attacks', 'https://lichess.org/training/CVDqa'),
('Discovered Attacks', 'https://lichess.org/training/bBRJf'),
('Discovered Attacks', 'https://lichess.org/training/oMYom'),
('Discovered Attacks', 'https://lichess.org/training/FKw3b'),

('Pins', 'https://lichess.org/training/zekfA'),
('Pins', 'https://lichess.org/training/elqPh'),
('Pins', 'https://lichess.org/training/zZ42x'),
('Pins', 'https://lichess.org/training/Qc1ks'),
('Pins', 'https://lichess.org/training/qHkft'),
('Pins', 'https://lichess.org/training/XqSb7'),
('Pins', 'https://lichess.org/training/pMxxr'),
('Pins', 'https://lichess.org/training/MoVic'),
('Pins', 'https://lichess.org/training/15qQb'),
('Pins', 'https://lichess.org/training/1ALnw'),
('Pins', 'https://lichess.org/training/kLgMh'),
('Pins', 'https://lichess.org/training/3KID3'),
('Pins', 'https://lichess.org/training/yWsxC'),
('Pins', 'https://lichess.org/training/um7dG'),
('Pins', 'https://lichess.org/training/aqyq2'),
('Pins', 'https://lichess.org/training/juhCX'),
('Pins', 'https://lichess.org/training/T9GfY'),
('Pins', 'https://lichess.org/training/bNQiH'),
('Pins', 'https://lichess.org/training/nXnpM'),
('Pins', 'https://lichess.org/training/syQLD'),

('Forks', 'https://lichess.org/training/2cpvV'),
('Forks', 'https://lichess.org/training/IDKxq'),
('Forks', 'https://lichess.org/training/p9BJl'),
('Forks', 'https://lichess.org/training/Vml8M'),
('Forks', 'https://lichess.org/training/WppiH'),
('Forks', 'https://lichess.org/training/y1PYP'),
('Forks', 'https://lichess.org/training/rbegm'),
('Forks', 'https://lichess.org/training/JSR02'),
('Forks', 'https://lichess.org/training/QUYMH'),
('Forks', 'https://lichess.org/training/5mpry'),
('Forks', 'https://lichess.org/training/t8a16'),
('Forks', 'https://lichess.org/training/CPgNC'),
('Forks', 'https://lichess.org/training/xNXpV'),
('Forks', 'https://lichess.org/training/yqUMx'),
('Forks', 'https://lichess.org/training/iMULS'),
('Forks', 'https://lichess.org/training/1R2Q5'),
('Forks', 'https://lichess.org/training/34dkq'),
('Forks', 'https://lichess.org/training/VqTmo'),
('Forks', 'https://lichess.org/training/YjXMk'),
('Forks', 'https://lichess.org/training/oGC3F'),

('Skewers', 'https://lichess.org/training/1bM2G'),
('Skewers', 'https://lichess.org/training/QZKJI'),
('Skewers', 'https://lichess.org/training/GVxtz'),
('Skewers', 'https://lichess.org/training/NgG5i'),
('Skewers', 'https://lichess.org/training/PqkRi'),
('Skewers', 'https://lichess.org/training/hlt2x'),
('Skewers', 'https://lichess.org/training/jS3pS'),
('Skewers', 'https://lichess.org/training/0ncTU'),
('Skewers', 'https://lichess.org/training/alXlk'),
('Skewers', 'https://lichess.org/training/MTkGp'),
('Skewers', 'https://lichess.org/training/FO8KF'),
('Skewers', 'https://lichess.org/training/CrFH0'),
('Skewers', 'https://lichess.org/training/pfed8'),
('Skewers', 'https://lichess.org/training/xQVPF'),
('Skewers', 'https://lichess.org/training/anL0N'),
('Skewers', 'https://lichess.org/training/qfqPR'),
('Skewers', 'https://lichess.org/training/LttYu'),
('Skewers', 'https://lichess.org/training/KkKXM'),
('Skewers', 'https://lichess.org/training/H5Nrk'),
('Skewers', 'https://lichess.org/training/EKq6x'),

('Rook and Pawn Endgame', 'https://lichess.org/training/eFWWS'),
('Rook and Pawn Endgame', 'https://lichess.org/training/4nZjQ'),
('Rook and Pawn Endgame', 'https://lichess.org/training/oZyc5'),
('Rook and Pawn Endgame', 'https://lichess.org/training/cgaXz'),
('Rook and Pawn Endgame', 'https://lichess.org/training/ibAmU'),
('Rook and Pawn Endgame', 'https://lichess.org/training/Ky4xC'),
('Rook and Pawn Endgame', 'https://lichess.org/training/4LDPr'),
('Rook and Pawn Endgame', 'https://lichess.org/training/hu6Kn'),
('Rook and Pawn Endgame', 'https://lichess.org/training/uOJGD'),
('Rook and Pawn Endgame', 'https://lichess.org/training/XyPu9');
('Rook and Pawn Endgame', 'https://lichess.org/training/C24wf'),
('Rook and Pawn Endgame', 'https://lichess.org/training/mSqKY'),
('Rook and Pawn Endgame', 'https://lichess.org/training/vgZ0H'),
('Rook and Pawn Endgame', 'https://lichess.org/training/n6y3I'),
('Rook and Pawn Endgame', 'https://lichess.org/training/Te6hB'),
('Rook and Pawn Endgame', 'https://lichess.org/training/I5o3j'),
('Rook and Pawn Endgame', 'https://lichess.org/training/ZYTqO'),
('Rook and Pawn Endgame', 'https://lichess.org/training/9MTjv'),
('Rook and Pawn Endgame', 'https://lichess.org/training/y5rch'),
('Rook and Pawn Endgame', 'https://lichess.org/training/v9z7O');