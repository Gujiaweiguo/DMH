import { defineComponent, h, ref, reactive, onBeforeUnmount, onMounted } from 'vue';
import * as LucideIcons from 'lucide-vue-next';
import { profileApi, type UserProfile } from '../services/profileApi';

export const UserProfileView = defineComponent({
  setup() {
    const currentUser = ref<UserProfile | null>(null);
    const loading = ref(false);
    const saving = ref(false);

    const isEditingProfile = ref(false);
    const isEditingPhone = ref(false);
    const isEditingEmail = ref(false);
    const isChangingPassword = ref(false);

    const profileForm = reactive({
      realName: '',
      username: '',
    });

    const phoneForm = reactive({
      newPhone: '',
      verifyCode: '',
    });

    const emailForm = reactive({
      newEmail: '',
      verifyCode: '',
    });

    const passwordForm = reactive({
      oldPassword: '',
      newPassword: '',
      confirmPassword: '',
    });

    const phoneCodeCountdown = ref(0);
    const emailCodeCountdown = ref(0);

    let phoneTimer: ReturnType<typeof setInterval> | null = null;
    let emailTimer: ReturnType<typeof setInterval> | null = null;

    const getErrorMessage = (error: unknown, fallback: string) => {
      if (error instanceof Error && error.message) {
        return error.message;
      }
      return fallback;
    };

    const startCountdown = (key: 'phone' | 'email') => {
      if (key === 'phone') {
        if (phoneTimer) {
          clearInterval(phoneTimer);
        }
        phoneCodeCountdown.value = 60;
        phoneTimer = setInterval(() => {
          phoneCodeCountdown.value -= 1;
          if (phoneCodeCountdown.value <= 0 && phoneTimer) {
            clearInterval(phoneTimer);
            phoneTimer = null;
          }
        }, 1000);
        return;
      }

      if (emailTimer) {
        clearInterval(emailTimer);
      }
      emailCodeCountdown.value = 60;
      emailTimer = setInterval(() => {
        emailCodeCountdown.value -= 1;
        if (emailCodeCountdown.value <= 0 && emailTimer) {
          clearInterval(emailTimer);
          emailTimer = null;
        }
      }, 1000);
    };

    const loadUserProfile = async () => {
      loading.value = true;
      try {
        currentUser.value = await profileApi.getUserInfo();
        profileForm.realName = currentUser.value.realName;
        profileForm.username = currentUser.value.username;
      } catch (error) {
        console.error('加载用户信息失败', error);
        alert(getErrorMessage(error, '加载用户信息失败'));
      } finally {
        loading.value = false;
      }
    };

    const saveProfile = async () => {
      if (!profileForm.realName.trim()) {
        alert('请输入姓名');
        return;
      }

      saving.value = true;
      try {
        const user = await profileApi.updateProfile(profileForm.realName.trim());
        currentUser.value = user;
        profileForm.realName = user.realName;
        isEditingProfile.value = false;
        alert('保存成功');
      } catch (error) {
        console.error('保存失败', error);
        alert(getErrorMessage(error, '保存失败，请重试'));
      } finally {
        saving.value = false;
      }
    };

    const sendPhoneCode = async () => {
      if (!phoneForm.newPhone) {
        alert('请输入手机号');
        return;
      }
      
      if (!/^1[3-9]\d{9}$/.test(phoneForm.newPhone)) {
        alert('请输入正确的手机号');
        return;
      }

      try {
        await profileApi.sendPhoneCode(phoneForm.newPhone);
        startCountdown('phone');
        alert('验证码已发送');
      } catch (error) {
        console.error('发送验证码失败', error);
        alert(getErrorMessage(error, '发送失败，请重试'));
      }
    };

    const bindPhone = async () => {
      if (!phoneForm.newPhone || !phoneForm.verifyCode) {
        alert('请填写完整信息');
        return;
      }

      saving.value = true;
      try {
        const user = await profileApi.bindPhone(phoneForm.newPhone, phoneForm.verifyCode);
        currentUser.value = user;
        isEditingPhone.value = false;
        phoneForm.newPhone = '';
        phoneForm.verifyCode = '';
        alert('绑定成功');
      } catch (error) {
        console.error('绑定失败', error);
        alert(getErrorMessage(error, '绑定失败，请检查验证码是否正确'));
      } finally {
        saving.value = false;
      }
    };

    const sendEmailCode = async () => {
      if (!emailForm.newEmail) {
        alert('请输入邮箱');
        return;
      }
      
      if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(emailForm.newEmail)) {
        alert('请输入正确的邮箱地址');
        return;
      }

      try {
        await profileApi.sendEmailCode(emailForm.newEmail);
        startCountdown('email');
        alert('验证码已发送到邮箱');
      } catch (error) {
        console.error('发送验证码失败', error);
        alert(getErrorMessage(error, '发送失败，请重试'));
      }
    };

    const bindEmail = async () => {
      if (!emailForm.newEmail || !emailForm.verifyCode) {
        alert('请填写完整信息');
        return;
      }

      saving.value = true;
      try {
        const user = await profileApi.bindEmail(emailForm.newEmail, emailForm.verifyCode);
        currentUser.value = user;
        isEditingEmail.value = false;
        emailForm.newEmail = '';
        emailForm.verifyCode = '';
        alert('绑定成功');
      } catch (error) {
        console.error('绑定失败', error);
        alert(getErrorMessage(error, '绑定失败，请检查验证码是否正确'));
      } finally {
        saving.value = false;
      }
    };

    const changePassword = async () => {
      if (!passwordForm.oldPassword || !passwordForm.newPassword || !passwordForm.confirmPassword) {
        alert('请填写完整信息');
        return;
      }
      
      if (passwordForm.newPassword.length < 6) {
        alert('新密码长度至少6位');
        return;
      }
      
      if (passwordForm.newPassword !== passwordForm.confirmPassword) {
        alert('两次输入的密码不一致');
        return;
      }

      saving.value = true;
      try {
        await profileApi.changePassword(passwordForm.oldPassword, passwordForm.newPassword);
        isChangingPassword.value = false;
        passwordForm.oldPassword = '';
        passwordForm.newPassword = '';
        passwordForm.confirmPassword = '';
        alert('密码修改成功，请重新登录');
      } catch (error) {
        console.error('修改密码失败', error);
        alert(getErrorMessage(error, '修改失败，请检查原密码是否正确'));
      } finally {
        saving.value = false;
      }
    };

    const getInputValue = (event: Event) => {
      return (event.target as HTMLInputElement).value;
    };

    const handleRealNameInput = (event: Event) => {
      profileForm.realName = getInputValue(event);
    };

    const handlePhoneInput = (event: Event) => {
      phoneForm.newPhone = getInputValue(event);
    };

    const handlePhoneCodeInput = (event: Event) => {
      phoneForm.verifyCode = getInputValue(event);
    };

    const handleEmailInput = (event: Event) => {
      emailForm.newEmail = getInputValue(event);
    };

    const handleEmailCodeInput = (event: Event) => {
      emailForm.verifyCode = getInputValue(event);
    };

    const handleOldPasswordInput = (event: Event) => {
      passwordForm.oldPassword = getInputValue(event);
    };

    const handleNewPasswordInput = (event: Event) => {
      passwordForm.newPassword = getInputValue(event);
    };

    const handleConfirmPasswordInput = (event: Event) => {
      passwordForm.confirmPassword = getInputValue(event);
    };

    const startEditProfile = () => {
      isEditingProfile.value = true;
    };

    const startEditPhone = () => {
      isEditingPhone.value = true;
    };

    const startEditEmail = () => {
      isEditingEmail.value = true;
    };

    const startChangePassword = () => {
      isChangingPassword.value = true;
    };

    onBeforeUnmount(() => {
      if (phoneTimer) {
        clearInterval(phoneTimer);
      }
      if (emailTimer) {
        clearInterval(emailTimer);
      }
    });

    onMounted(() => {
      void loadUserProfile();
    });

    return () => {
      if (loading.value) {
        return h('div', { class: 'flex items-center justify-center h-96' }, [
          h('div', { class: 'text-center' }, [
            h('div', { class: 'animate-spin rounded-full h-12 w-12 border-4 border-indigo-600 border-t-transparent mx-auto' }),
            h('p', { class: 'mt-4 text-slate-600' }, '加载中...')
          ])
        ]);
      }

      if (!currentUser.value) {
        return h('div', { class: 'text-center py-20 text-slate-400' }, '加载用户信息失败');
      }

      return h('div', { class: 'max-w-4xl mx-auto space-y-6 animate-in fade-in' }, [
        // 页面标题
        h('div', [h('h2', { class: 'text-4xl font-black text-slate-900' }, '个人设置'), h('p', { class: 'text-slate-400 mt-2' }, '管理你的个人资料和账号安全')]),

        // 用户头像和基本信息
        h('div', { class: 'bg-white rounded-3xl p-8 border border-slate-100' }, [
          h('div', { class: 'flex items-center gap-6' }, [
            h('img', {
              src: currentUser.value.avatar,
              class: 'w-24 h-24 rounded-3xl border-4 border-white shadow-xl'
            }),
            h('div', { class: 'flex-1' }, [
              h('h3', { class: 'text-2xl font-black text-slate-900' }, currentUser.value.realName),
              h('p', { class: 'text-slate-500 mt-1' }, `@${currentUser.value.username}`),
              h('div', { class: 'flex gap-2 mt-3' },
                currentUser.value.roles.map((role: string) => {
                  const colors: Record<string, string> = {
                    platform_admin: 'bg-purple-100 text-purple-700 border-purple-200',
                    participant: 'bg-green-100 text-green-700 border-green-200',
                  };
                  const names: Record<string, string> = {
                    platform_admin: '平台管理员',
                    participant: '活动参与者',
                  };
                  return h('span', {
                    class: `px-3 py-1 rounded-full text-xs font-bold border ${colors[role] || 'bg-gray-100 text-gray-700'}`
                  }, names[role] || role);
                })
              )
            ])
          ])
        ]),

        // 个人资料卡片
        h('div', { class: 'bg-white rounded-3xl p-8 border border-slate-100' }, [
          h('div', { class: 'flex items-center justify-between mb-6' }, [
            h('h3', { class: 'text-xl font-black text-slate-900 flex items-center gap-2' }, [
              h(LucideIcons.User, { size: 20, class: 'text-indigo-600' }),
              '个人资料'
            ]),
            !isEditingProfile.value && h('button', {
              onClick: startEditProfile,
              class: 'text-sm font-bold text-indigo-600 hover:text-indigo-500 flex items-center gap-1'
            }, [h(LucideIcons.Edit, { size: 16 }), '编辑'])
          ]),

          isEditingProfile.value
            ? h('div', { class: 'space-y-4' }, [
                h('div', [
                  h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '真实姓名'),
                  h('input', {
                    type: 'text',
                    value: profileForm.realName,
                    onInput: handleRealNameInput,
                    placeholder: '请输入真实姓名',
                    class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors'
                  })
                ]),
                h('div', [
                  h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '用户名'),
                  h('input', {
                    type: 'text',
                    value: profileForm.username,
                    disabled: true,
                    class: 'w-full px-4 py-3 rounded-xl border border-slate-200 bg-slate-50 text-slate-500'
                  }),
                  h('p', { class: 'mt-1 text-xs text-slate-400' }, '用户名不可修改')
                ]),
                h('div', { class: 'flex gap-3 pt-2' }, [
                  h('button', {
                    onClick: () => {
                      isEditingProfile.value = false;
                      profileForm.realName = currentUser.value.realName;
                    },
                    class: 'flex-1 px-6 py-3 rounded-xl border border-slate-200 font-bold hover:bg-slate-50 transition-colors'
                  }, '取消'),
                  h('button', {
                    onClick: saveProfile,
                    disabled: saving.value,
                    class: 'flex-1 px-6 py-3 rounded-xl bg-indigo-600 text-white font-bold hover:bg-indigo-500 transition-colors disabled:opacity-50'
                  }, saving.value ? '保存中...' : '保存')
                ])
              ])
            : h('div', { class: 'space-y-3' }, [
                h('div', { class: 'flex justify-between py-3 border-b border-slate-100' }, [
                  h('span', { class: 'text-slate-600' }, '真实姓名'),
                  h('span', { class: 'font-bold text-slate-900' }, currentUser.value.realName)
                ]),
                h('div', { class: 'flex justify-between py-3' }, [
                  h('span', { class: 'text-slate-600' }, '用户名'),
                  h('span', { class: 'font-bold text-slate-900' }, currentUser.value.username)
                ])
              ])
        ]),

        // 手机号卡片
        h('div', { class: 'bg-white rounded-3xl p-8 border border-slate-100' }, [
          h('div', { class: 'flex items-center justify-between mb-6' }, [
            h('h3', { class: 'text-xl font-black text-slate-900 flex items-center gap-2' }, [
              h(LucideIcons.Smartphone, { size: 20, class: 'text-emerald-600' }),
              '手机号'
            ]),
            !isEditingPhone.value && h('button', {
              onClick: startEditPhone,
              class: 'text-sm font-bold text-indigo-600 hover:text-indigo-500 flex items-center gap-1'
            }, [h(LucideIcons.Edit, { size: 16 }), currentUser.value.phone ? '换绑' : '绑定'])
          ]),

          isEditingPhone.value
            ? h('div', { class: 'space-y-4' }, [
                h('div', { class: 'bg-blue-50 border border-blue-200 rounded-xl p-4' }, [
                  h('div', { class: 'flex items-start gap-3' }, [
                    h(LucideIcons.Info, { size: 20, class: 'text-blue-600 flex-shrink-0 mt-0.5' }),
                    h('div', { class: 'text-sm text-blue-800' }, [
                      h('p', { class: 'font-bold' }, currentUser.value.phone ? '换绑手机号' : '绑定手机号'),
                      h('p', { class: 'mt-1' }, '验证码将发送到新手机号，请注意查收')
                    ])
                  ])
                ]),
                h('div', [
                  h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '新手机号'),
                  h('input', {
                    type: 'tel',
                    value: phoneForm.newPhone,
                    onInput: handlePhoneInput,
                    placeholder: '请输入新手机号',
                    class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors'
                  })
                ]),
                h('div', [
                  h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '验证码'),
                  h('div', { class: 'flex gap-2' }, [
                    h('input', {
                      type: 'text',
                      value: phoneForm.verifyCode,
                      onInput: handlePhoneCodeInput,
                      placeholder: '请输入验证码',
                      class: 'flex-1 px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors'
                    }),
                    h('button', {
                      onClick: sendPhoneCode,
                      disabled: phoneCodeCountdown.value > 0,
                      class: 'px-6 py-3 rounded-xl bg-emerald-600 text-white font-bold hover:bg-emerald-500 transition-colors disabled:opacity-50 disabled:cursor-not-allowed whitespace-nowrap'
                    }, phoneCodeCountdown.value > 0 ? `${phoneCodeCountdown.value}s` : '获取验证码')
                  ])
                ]),
                h('div', { class: 'flex gap-3 pt-2' }, [
                  h('button', {
                    onClick: () => {
                      isEditingPhone.value = false;
                      phoneForm.newPhone = '';
                      phoneForm.verifyCode = '';
                    },
                    class: 'flex-1 px-6 py-3 rounded-xl border border-slate-200 font-bold hover:bg-slate-50 transition-colors'
                  }, '取消'),
                  h('button', {
                    onClick: bindPhone,
                    disabled: saving.value,
                    class: 'flex-1 px-6 py-3 rounded-xl bg-emerald-600 text-white font-bold hover:bg-emerald-500 transition-colors disabled:opacity-50'
                  }, saving.value ? '绑定中...' : '确认绑定')
                ])
              ])
            : h('div', { class: 'flex justify-between items-center py-3' }, [
                h('span', { class: 'text-slate-600' }, '已绑定手机'),
                h('span', { class: 'font-bold text-slate-900' }, currentUser.value.phone || '未绑定')
              ])
        ]),

        // 邮箱卡片
        h('div', { class: 'bg-white rounded-3xl p-8 border border-slate-100' }, [
          h('div', { class: 'flex items-center justify-between mb-6' }, [
            h('h3', { class: 'text-xl font-black text-slate-900 flex items-center gap-2' }, [
              h(LucideIcons.Mail, { size: 20, class: 'text-purple-600' }),
              '邮箱'
            ]),
            !isEditingEmail.value && h('button', {
              onClick: startEditEmail,
              class: 'text-sm font-bold text-indigo-600 hover:text-indigo-500 flex items-center gap-1'
            }, [h(LucideIcons.Edit, { size: 16 }), currentUser.value.email ? '换绑' : '绑定'])
          ]),

          isEditingEmail.value
            ? h('div', { class: 'space-y-4' }, [
                h('div', { class: 'bg-purple-50 border border-purple-200 rounded-xl p-4' }, [
                  h('div', { class: 'flex items-start gap-3' }, [
                    h(LucideIcons.Info, { size: 20, class: 'text-purple-600 flex-shrink-0 mt-0.5' }),
                    h('div', { class: 'text-sm text-purple-800' }, [
                      h('p', { class: 'font-bold' }, currentUser.value.email ? '换绑邮箱' : '绑定邮箱'),
                      h('p', { class: 'mt-1' }, '验证码将发送到新邮箱，请注意查收')
                    ])
                  ])
                ]),
                h('div', [
                  h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '新邮箱'),
                  h('input', {
                    type: 'email',
                    value: emailForm.newEmail,
                    onInput: handleEmailInput,
                    placeholder: '请输入新邮箱',
                    class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors'
                  })
                ]),
                h('div', [
                  h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '验证码'),
                  h('div', { class: 'flex gap-2' }, [
                    h('input', {
                      type: 'text',
                      value: emailForm.verifyCode,
                      onInput: handleEmailCodeInput,
                      placeholder: '请输入验证码',
                      class: 'flex-1 px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors'
                    }),
                    h('button', {
                      onClick: sendEmailCode,
                      disabled: emailCodeCountdown.value > 0,
                      class: 'px-6 py-3 rounded-xl bg-purple-600 text-white font-bold hover:bg-purple-500 transition-colors disabled:opacity-50 disabled:cursor-not-allowed whitespace-nowrap'
                    }, emailCodeCountdown.value > 0 ? `${emailCodeCountdown.value}s` : '获取验证码')
                  ])
                ]),
                h('div', { class: 'flex gap-3 pt-2' }, [
                  h('button', {
                    onClick: () => {
                      isEditingEmail.value = false;
                      emailForm.newEmail = '';
                      emailForm.verifyCode = '';
                    },
                    class: 'flex-1 px-6 py-3 rounded-xl border border-slate-200 font-bold hover:bg-slate-50 transition-colors'
                  }, '取消'),
                  h('button', {
                    onClick: bindEmail,
                    disabled: saving.value,
                    class: 'flex-1 px-6 py-3 rounded-xl bg-purple-600 text-white font-bold hover:bg-purple-500 transition-colors disabled:opacity-50'
                  }, saving.value ? '绑定中...' : '确认绑定')
                ])
              ])
            : h('div', { class: 'flex justify-between items-center py-3' }, [
                h('span', { class: 'text-slate-600' }, '已绑定邮箱'),
                h('span', { class: 'font-bold text-slate-900' }, currentUser.value.email || '未绑定')
              ])
        ]),

        // 密码修改卡片
        h('div', { class: 'bg-white rounded-3xl p-8 border border-slate-100' }, [
          h('div', { class: 'flex items-center justify-between mb-6' }, [
            h('h3', { class: 'text-xl font-black text-slate-900 flex items-center gap-2' }, [
              h(LucideIcons.Lock, { size: 20, class: 'text-amber-600' }),
              '登录密码'
            ]),
            !isChangingPassword.value && h('button', {
              onClick: startChangePassword,
              class: 'text-sm font-bold text-indigo-600 hover:text-indigo-500 flex items-center gap-1'
            }, [h(LucideIcons.Key, { size: 16 }), '修改密码'])
          ]),

          isChangingPassword.value
            ? h('div', { class: 'space-y-4' }, [
                h('div', { class: 'bg-amber-50 border border-amber-200 rounded-xl p-4' }, [
                  h('div', { class: 'flex items-start gap-3' }, [
                    h(LucideIcons.AlertTriangle, { size: 20, class: 'text-amber-600 flex-shrink-0 mt-0.5' }),
                    h('div', { class: 'text-sm text-amber-800' }, [
                      h('p', { class: 'font-bold' }, '修改登录密码'),
                      h('p', { class: 'mt-1' }, '密码修改成功后需要重新登录')
                    ])
                  ])
                ]),
                h('div', [
                  h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '原密码'),
                  h('input', {
                    type: 'password',
                    value: passwordForm.oldPassword,
                    onInput: handleOldPasswordInput,
                    placeholder: '请输入原密码',
                    class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors'
                  })
                ]),
                h('div', [
                  h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '新密码'),
                  h('input', {
                    type: 'password',
                    value: passwordForm.newPassword,
                    onInput: handleNewPasswordInput,
                    placeholder: '请输入新密码（至少6位）',
                    class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors'
                  })
                ]),
                h('div', [
                  h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '确认新密码'),
                  h('input', {
                    type: 'password',
                    value: passwordForm.confirmPassword,
                    onInput: handleConfirmPasswordInput,
                    placeholder: '请再次输入新密码',
                    class: 'w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:border-indigo-500 transition-colors'
                  })
                ]),
                h('div', { class: 'flex gap-3 pt-2' }, [
                  h('button', {
                    onClick: () => {
                      isChangingPassword.value = false;
                      passwordForm.oldPassword = '';
                      passwordForm.newPassword = '';
                      passwordForm.confirmPassword = '';
                    },
                    class: 'flex-1 px-6 py-3 rounded-xl border border-slate-200 font-bold hover:bg-slate-50 transition-colors'
                  }, '取消'),
                  h('button', {
                    onClick: changePassword,
                    disabled: saving.value,
                    class: 'flex-1 px-6 py-3 rounded-xl bg-amber-600 text-white font-bold hover:bg-amber-500 transition-colors disabled:opacity-50'
                  }, saving.value ? '修改中...' : '确认修改')
                ])
              ])
            : h('div', { class: 'flex justify-between items-center py-3' }, [
                h('span', { class: 'text-slate-600' }, '登录密码'),
                h('span', { class: 'font-bold text-slate-900' }, '••••••••')
              ])
        ]),

        // 账号信息
        h('div', { class: 'bg-white rounded-3xl p-8 border border-slate-100' }, [
          h('h3', { class: 'text-xl font-black text-slate-900 flex items-center gap-2 mb-6' }, [
            h(LucideIcons.Info, { size: 20, class: 'text-slate-600' }),
            '账号信息'
          ]),
          h('div', { class: 'space-y-3' }, [
            h('div', { class: 'flex justify-between py-3 border-b border-slate-100' }, [
              h('span', { class: 'text-slate-600' }, '用户ID'),
              h('span', { class: 'font-bold text-slate-900' }, currentUser.value.id)
            ]),
            h('div', { class: 'flex justify-between py-3' }, [
              h('span', { class: 'text-slate-600' }, '注册时间'),
              h('span', { class: 'font-bold text-slate-900' }, currentUser.value.createdAt)
            ])
          ])
        ])
      ]);
    };
  }
});
