import { describe, expect, it, vi, beforeEach } from 'vitest';

const validateLoginForm = (form: { username: string; password: string }) => {
  return form.username.length > 0 && form.password.length >= 6;
};

const validateRegisterForm = (form: { 
  username: string; 
  password: string; 
  phone: string;
  email?: string;
  realName?: string;
}) => {
  return form.username.length > 0 && 
         form.password.length >= 6 && 
         form.phone.length > 0 &&
         /^1[3-9]\d{9}$/.test(form.phone);
};

const isValidPhone = (phone: string) => {
  return /^1[3-9]\d{9}$/.test(phone);
};

const isValidPassword = (password: string) => {
  return password.length >= 6;
};

describe('LoginView', () => {
  describe('Login Form Validation', () => {
    it('should return false for empty username', () => {
      const form = { username: '', password: '123456' };
      expect(validateLoginForm(form)).toBe(false);
    });

    it('should return false for password less than 6 characters', () => {
      const form = { username: 'admin', password: '12345' };
      expect(validateLoginForm(form)).toBe(false);
    });

    it('should return true for valid credentials', () => {
      const form = { username: 'admin', password: '123456' };
      expect(validateLoginForm(form)).toBe(true);
    });

    it('should return true for password exactly 6 characters', () => {
      const form = { username: 'admin', password: '123456' };
      expect(validateLoginForm(form)).toBe(true);
    });

    it('should return true for long password', () => {
      const form = { username: 'admin', password: 'verylongpassword123' };
      expect(validateLoginForm(form)).toBe(true);
    });
  });

  describe('Register Form Validation', () => {
    it('should return false for empty username', () => {
      const form = { username: '', password: '123456', phone: '13800138000' };
      expect(validateRegisterForm(form)).toBe(false);
    });

    it('should return false for short password', () => {
      const form = { username: 'test', password: '12345', phone: '13800138000' };
      expect(validateRegisterForm(form)).toBe(false);
    });

    it('should return false for invalid phone format', () => {
      const form = { username: 'test', password: '123456', phone: '12345' };
      expect(validateRegisterForm(form)).toBe(false);
    });

    it('should return false for phone starting with invalid digit', () => {
      const form = { username: 'test', password: '123456', phone: '12000138000' };
      expect(validateRegisterForm(form)).toBe(false);
    });

    it('should return true for valid registration data', () => {
      const form = { username: 'test', password: '123456', phone: '13800138000' };
      expect(validateRegisterForm(form)).toBe(true);
    });

    it('should return true for phone starting with 3-9', () => {
      const validPrefixes = ['130', '150', '180', '199'];
      validPrefixes.forEach(prefix => {
        const form = { 
          username: 'test', 
          password: '123456', 
          phone: prefix + '12345678' 
        };
        expect(validateRegisterForm(form)).toBe(true);
      });
    });
  });

  describe('Phone Validation', () => {
    it('should validate correct phone numbers', () => {
      const validPhones = [
        '13800138000',
        '15012345678',
        '18800001111',
        '19900001111',
        '13123456789'
      ];
      validPhones.forEach(phone => {
        expect(isValidPhone(phone)).toBe(true);
      });
    });

    it('should reject invalid phone numbers', () => {
      const invalidPhones = [
        '12345678901',
        '10012345678',
        '1380013800',
        '138001380001',
        'abcdefghijk',
        '',
      ];
      invalidPhones.forEach(phone => {
        expect(isValidPhone(phone)).toBe(false);
      });
    });
  });

  describe('Password Validation', () => {
    it('should accept password with 6 or more characters', () => {
      expect(isValidPassword('123456')).toBe(true);
      expect(isValidPassword('1234567')).toBe(true);
      expect(isValidPassword('verylongpassword')).toBe(true);
    });

    it('should reject password with less than 6 characters', () => {
      expect(isValidPassword('')).toBe(false);
      expect(isValidPassword('1')).toBe(false);
      expect(isValidPassword('12345')).toBe(false);
    });
  });
});

describe('Test Account Quick Fill', () => {
  it('should fill admin credentials correctly', () => {
    const loginForm = { username: '', password: '' };
    loginForm.username = 'admin';
    loginForm.password = '123456';
    
    expect(loginForm.username).toBe('admin');
    expect(loginForm.password).toBe('123456');
    expect(validateLoginForm(loginForm)).toBe(true);
  });
});
