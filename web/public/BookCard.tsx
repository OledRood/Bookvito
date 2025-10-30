import React, { useState } from 'react';
import { Card, CardActionArea, CardContent, CardMedia, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import resolveImageUrl from '../src/utils/imageUrl';

interface BookCardProps {
  id: string | number;
  imageUrl: string;
  title: string;
  author: string;
}

const BookCard: React.FC<BookCardProps> = ({ id, imageUrl, title, author }) => {
  const [hasError, setHasError] = useState(false);
  const fallback = '/images/default-book.png';
  const imgSrc = !imageUrl || hasError ? fallback : resolveImageUrl(imageUrl);

  return (
    <Card
      elevation={1}
      sx={{
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        borderRadius: 'var(--radii-medium, 12px)',
        transition: `transform var(--motion-short,200ms) ease, box-shadow var(--motion-short,200ms) ease`,
        '&:hover': { transform: 'translateY(-6px)', boxShadow: 'var(--shadow-dp2)' },
      }}
    >
      <CardActionArea
        component={RouterLink}
        to={`/book/${id}`}
        sx={{ display: 'flex', flexDirection: 'column', height: '100%', justifyContent: 'flex-start', alignItems: 'flex-start' }}
      >
        <CardMedia
          component="img"
          sx={{
            // slightly taller image on very small screens but keep reasonable content size
            height: { xs: 210, sm: 190 },
            objectFit: 'cover',
            borderTopLeftRadius: 'var(--radii-medium, 12px)',
            borderTopRightRadius: 'var(--radii-medium, 12px)',
          }}
          image={imgSrc}
          alt={title}
          onError={(e: React.SyntheticEvent<HTMLImageElement, Event>) => {
            // ensure we only handle once and set fallback
            const target = e.currentTarget as HTMLImageElement;
            if (target.src && !hasError) {
              target.onerror = null;
              target.src = fallback;
              setHasError(true);
            }
          }}
        />
        <CardContent sx={{ flexGrow: 1, textAlign: 'left', display: 'flex', flexDirection: 'column', alignItems: 'flex-start', mt: { xs: '20px', md: 0 } }}>
          <Typography
            gutterBottom
            variant="h6"
            component="div"
            sx={{
              // Ограничение названия двумя строками
              display: '-webkit-box',
              WebkitLineClamp: 2,
              WebkitBoxOrient: 'vertical',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              minHeight: { xs: '2.2em', md: '3.2em' }, // smaller on mobile to avoid large blocks
              fontSize: { xs: '1rem', md: 'unset' },
              width: '100%',
              textAlign: 'left',
              fontWeight: 600,
            }}
          >
            {title}
          </Typography>
          <Typography variant="body2" color="text.secondary" align="left" sx={{ width: '100%', textAlign: 'left' }}>
            {author}
          </Typography>
        </CardContent>
      </CardActionArea>
    </Card>
  );
};

export default BookCard;