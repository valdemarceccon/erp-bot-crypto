import React from 'react';
import { useForm } from 'react-hook-form';
import { Button, TextField, Grid, Paper, Typography } from '@mui/material';

interface LoginProps {
  onLogin: (email: string, password: string) => void;
}

interface LoginFormData {
  email: string;
  password: string;
}

const LoginForm: React.FC<LoginProps> = ({ onLogin }) => {
  const { register, handleSubmit, formState: { errors } } = useForm<LoginFormData>();

  const onSubmit = (data: LoginFormData) => {
    onLogin(data.email, data.password);
  };

  return (
    <Grid
      container
      justifyContent="center">
      <Grid item xs={12} sm={8} md={6} lg={4}>
        <Paper elevation={3} sx={{ padding: 2 }}>
          <Typography variant="h5" gutterBottom>
            Login
          </Typography>
          <form onSubmit={handleSubmit(onSubmit)}>
            <TextField
              fullWidth
              margin="normal"
              type="email"
              label="Email Address"
              {...register('email', { required: 'Email is required' })}
              error={!!errors.email}
              helperText={errors.email?.message}
            />
            <TextField
              fullWidth
              margin="normal"
              type="password"
              label="Password"
              {...register('password', { required: 'Password is required' })}
              error={!!errors.password}
              helperText={errors.password?.message}
            />
            <Button
              fullWidth
              type="submit"
              variant="contained"
              color="primary"
              sx={{ mt: 2 }}
            >
              Login
            </Button>
          </form>
        </Paper>
      </Grid>
    </Grid>
  );
};

export default LoginForm;
