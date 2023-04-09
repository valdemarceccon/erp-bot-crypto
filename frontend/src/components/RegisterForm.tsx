import React from 'react';
import { useForm } from 'react-hook-form';
import { Button, TextField, Grid, Paper, Typography } from '@mui/material';

interface RegisterProps {
  onRegister: (email: string, name: string, password: string) => void;
}

interface FormInputs {
  email: string;
  name: string;
  password: string;
  confirmPassword: string;
}

const RegisterForm: React.FC<RegisterProps> = ({ onRegister }) => {
  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
  } = useForm<FormInputs>();

  const password = watch('password', '');

  const onSubmit = (data: FormInputs) => {
    if (data.password === data.confirmPassword) {
      onRegister(data.email, data.name, data.password);
    }
  };

  return (
    <Grid container justifyContent="center">
      <Grid item xs={12} sm={8} md={6} lg={4}>
        <Paper elevation={3} sx={{ padding: 2 }}>
          <Typography variant="h5" gutterBottom>
            Register
          </Typography>
          <form onSubmit={handleSubmit(onSubmit)}>
            <TextField
              fullWidth
              margin="normal"
              type="email"
              label="Email Address"
              {...register('email', { required: 'Email is required.' })}
              error={!!errors.email}
              helperText={errors.email?.message}
            />
            <TextField
              fullWidth
              margin="normal"
              type="text"
              label="Name"
              {...register('name', { required: 'Name is required.' })}
              error={!!errors.name}
              helperText={errors.name?.message}
            />
            <TextField
              fullWidth
              margin="normal"
              type="password"
              label="Password"
              {...register('password', { required: 'Password is required.' })}
              error={!!errors.password}
              helperText={errors.password?.message}
            />
            <TextField
              fullWidth
              margin="normal"
              type="password"
              label="Confirm Password"
              {...register('confirmPassword', {
                required: 'Confirm Password is required.',
                validate: (value) =>
                  value === password || 'Passwords do not match.',
              })}
              error={!!errors.confirmPassword}
              helperText={errors.confirmPassword?.message}
            />
            <Button
              fullWidth
              type="submit"
              variant="contained"
              color="primary"
              sx={{ mt: 2 }}
            >
              Register
            </Button>
          </form>
        </Paper>
      </Grid>
    </Grid>
  );
};

export default RegisterForm;
