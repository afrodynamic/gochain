import { Card, CardContent, Stack, Typography } from '@mui/material';

export const MetricCard = ({
  label,
  value,
  helper,
}: {
  label: string;
  value: string | number;
  helper?: string;
}) => {
  return (
    <Card
      elevation={6}
      className="max-w-lg w-full border"
      sx={{ borderColor: 'primary.main' }}
    >
      <CardContent>
        <Stack spacing={1}>
          <Typography variant="body2" fontWeight="bold">
            {label}
          </Typography>

          <Typography
            variant="h4"
            fontWeight={700}
            sx={{ wordBreak: 'break-word' }}
          >
            {value}
          </Typography>

          {helper && <Typography variant="caption">{helper}</Typography>}
        </Stack>
      </CardContent>
    </Card>
  );
};
