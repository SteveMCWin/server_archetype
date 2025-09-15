DELETE from quotes;
DELETE FROM sqlite_sequence WHERE name='quotes';

-- CHEST
INSERT INTO quotes (quote, source, len) VALUES
('Be yourself; everyone else is already taken.', 'Oscar Wilde', 0),
('I''m selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hard to handle. But if you can''t handle me at my worst, then you sure as hell don''t deserve me at my best.', 'Marilyn Monroe', 1),
('So many books, so little time.', 'Frank Zappa', 0),
('Two things are infinite: the universe and human stupidity; and I''m not sure about the universe.', 'Albert Einstein', 0),
('You know you''re in love when you can''t fall asleep because reality is finally better than your dreams.', 'Dr. Seuss', 0),
('Darkness cannot drive out darkness: only light can do that. Hate cannot drive out hate: only love can do that.', 'Martin Luther King Jr., A Testament of Hope: The Essential Writing', 0),
('It was the best of times, it was the worst of times, it was the age of wisdom, it was the age of foolishness, it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of Darkness, it was the spring of hope, it was the winter of despair.', 'Charles Dickens, A Tale of Two Cities', 1),
('Beware; for I am fearless, and therefore powerful.', 'Mary Shelley, Frankenstein', 0),
('I wanted you to see what real courage is, instead of getting the idea that courage is a man with a gun in his hand. It''s when you know you''re licked before you begin but you begin anyway and you see it through no matter what. You rarely win, but sometimes you do.', 'Harper Lee, To Kill a Mockingbird', 1),
('The same substance composes us — the tree overhead, the stone beneath us, the bird, the beast, the star — we are all one, all moving to the same end.', 'P.L. Travers, Mary Poppins', 1);
