import { Container } from '@mui/material';
import clsx from 'clsx';
import { ReactNode } from 'react';

interface PageSectionProps {
  id: string;
  full?: boolean;
  className?: string;
  children: ReactNode;
}

export function PageSection({
  id,
  full = false,
  className,
  children,
}: PageSectionProps) {
  return (
    <section
      id={id}
      className={clsx(
        'w-full',
        full && 'min-h-svh',
        'py-8 md:py-12 px-4 sm:px-6',
        className ?? ''
      )}
    >
      <Container
        maxWidth="xl"
        className="flex flex-grow flex-col items-center h-full"
      >
        {children}
      </Container>
    </section>
  );
}
