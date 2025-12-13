import z from 'zod';

export const minLengthErrorMessage =
	'Password must be at least 8 characters long.';

export const maxLengthErrorMessage =
	'Password must be no more than 20 characters long.';

export const uppercaseErrorMessage =
	'Password must contain at least one uppercase letter (A–Z).';

export const lowercaseErrorMessage =
	'Password must contain at least one lowercase letter (a–z).';

export const numberErrorMessage =
	'Password must contain at least one number (0–9).';

export const specialCharacterErrorMessage =
	'Password must contain at least one special character (!@#$%^&*).';

export const passwordMismatchErrorMessage = 'Passwords do not match.';

export const passwordSchema = z
	.string()
	.min(8, { error: minLengthErrorMessage })
	.max(20, { error: maxLengthErrorMessage })
	.refine((password) => /[A-Z]/.test(password), {
		error: uppercaseErrorMessage,
	})
	.refine((password) => /[a-z]/.test(password), {
		error: lowercaseErrorMessage,
	})
	.refine((password) => /[0-9]/.test(password), { error: numberErrorMessage })
	.refine((password) => /[!@#$%^&*]/.test(password), {
		error: specialCharacterErrorMessage,
	});

export const updatePasswordSchema = z
	.object({
		currentPassword: z.string(),
		password: passwordSchema,
		confirmPassword: z.string(),
	})
	.refine((data) => data.password === data.confirmPassword, {
		error: passwordMismatchErrorMessage,
		path: ['confirmPassword'],
	});
