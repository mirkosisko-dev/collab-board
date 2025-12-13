import { useMutation } from '@tanstack/react-query';
import userApi from '../service';
import type { LoginPayload } from '../types';

const useLoginUser = () =>
	useMutation({
		mutationFn: (payload: LoginPayload) => userApi.login(payload),
	});

export default useLoginUser;
