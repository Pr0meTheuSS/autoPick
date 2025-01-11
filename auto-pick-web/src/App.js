import React, { useState } from 'react';
import { Container, TextField, Button, Grid, Card, CardContent, CardMedia, Typography, CircularProgress } from '@mui/material';
import axios from 'axios';

function App() {
  const [searchString, setSearchString] = useState('');
  const [model, setModel] = useState('');
  const [parts, setParts] = useState([]);
  const [loading, setLoading] = useState(false);

  const fetchParts = async () => {
    setLoading(true);
    try {
      const response = await axios.get(`http://localhost:8000/drom/${searchString}/${model}`);
      setParts(response.data);
    } catch (error) {
      console.error('Error fetching parts:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container maxWidth="md" style={{ marginTop: '40px' }}>
      <Typography variant="h4" align="center" gutterBottom>
        Агрегатор запчастей
      </Typography>
      <Grid container spacing={2} justifyContent="center">
        <Grid item xs={12} sm={5}>
          <TextField
            fullWidth
            label="Поиск запчастей"
            variant="outlined"
            value={searchString}
            onChange={(e) => setSearchString(e.target.value)}
          />
        </Grid>
        <Grid item xs={12} sm={5}>
          <TextField
            fullWidth
            label="Модель"
            variant="outlined"
            value={model}
            onChange={(e) => setModel(e.target.value)}
          />
        </Grid>
        <Grid item xs={12} sm={2}>
          <Button fullWidth variant="contained" color="primary" onClick={fetchParts}>
            Найти
          </Button>
        </Grid>
      </Grid>

      {loading ? (
        <Grid container justifyContent="center" style={{ marginTop: '20px' }}>
          <CircularProgress />
        </Grid>
      ) : (
        <Grid container spacing={2} style={{ marginTop: '20px' }}>
          {parts.length > 0 ? (
            parts.map((part, index) => (
              <Grid item xs={12} sm={6} md={4} key={index}>
                <Card>
                  <CardMedia
                    component="img"
                    height="140"
                    image={part.image_url}
                    alt={part.title}
                  />
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      {part.title}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      Бренд: {part.brand}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      Цена: {part.price}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      Местоположение: {part.location}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      Дата: {part.date}
                    </Typography>
                    <Button size="small" color="primary" href={part.link} target="_blank">
                      Перейти
                    </Button>
                  </CardContent>
                </Card>
              </Grid>
            ))
          ) : (
            <Typography variant="body1" align="center" style={{ marginTop: '10px' }}>
              Нет данных для отображения
            </Typography>
          )}
        </Grid>
      )}
    </Container>
  );
}

export default App;
