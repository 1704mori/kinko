import { twMerge } from 'tailwind-merge';
import clsx from 'clsx';

export function cn(...classes: any[]) {
  return clsx(twMerge(...classes));
}

export async function to<T = any>(promise: Promise<any>): Promise<[T | null, Error | null]> {
  try {
    const res = await promise;
    return [res, null];
  } catch (error) {
    return [null, error as Error];
  }
}