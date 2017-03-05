# otp

[![Join the chat at https://gitter.im/adminfromhell-otp/Lobby](https://badges.gitter.im/adminfromhell-otp/Lobby.svg)](https://gitter.im/adminfromhell-otp/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
<!-- *(Code Health and Pipleine Status at bottem of README)*   -->

CLI based TOTP tool  

### Initalize OTP
run `otp init`
you will be asked to supply a password that will be used to encrypt the database that holds the OTP secrets.

### Add an account to OTP
run `otp add --nickname "some nickname"`
Enter the information that is prompted on the screen.

### Get OTP codes
run `otp` without and arguments
enter your password and it will display a list of otp codes as well as thier timer.

### List accounts
run `otp list`
enter your password and you will be shown a list of what accounts you have stored.

### Notes
OTP stores its database inside ~/.otp
vendoring is done with Glide

### The TODO List
- [ ] Add delete function to remove old/bad keys without
- [ ] Add way to update OTP secrets without needing to delete entry
- [ ] Finish Export functionality
- [ ] Add password zip encryption to exports
- [ ] Add way to specify password without entering on CLI
- [ ] Sync database with WebDAV (ownCloud/NextCloud & etc.)
- [ ] Add HOTP functionality
- [ ] Get project above 50% unit testing
- [ ] Get project into Concourse CI
- [ ] Provide up-to date RPM, DEB, AUR packages if needed

### Credits
@erasche for his project that gave me some insperation when starting this one.













<!-- ---
### Code Insight  
###### Master Branch:  
[![Code Climate](https://codeclimate.com/github/adminfromhell/otp/badges/gpa.svg)](https://codeclimate.com/github/adminfromhell/otp) [![Issue Count](https://codeclimate.com/github/adminfromhell/otp/badges/issue_count.svg)](https://codeclimate.com/github/adminfromhell/otp)  

| Pipleine Job | Status |
|:------------:|:------:|
| Tests | [![Go Tests](https://ci.mythic.tech/api/v1/teams/Mythic%20Tech/pipelines/otp/jobs/test-master/badge)](https://ci.mythic.tech/teams/Mythic%20Tech/pipelines/otp) |
| Build | [![Go Tests](https://ci.mythic.tech/api/v1/teams/Mythic%20Tech/pipelines/otp/jobs/build-master/badge)](https://ci.mythic.tech/teams/Mythic%20Tech/pipelines/otp) ||

###### Develop Branch:  

| Pipleine Job | Status |
|:------------:|:------:|
| Tests | [![Go Tests](https://ci.mythic.tech/api/v1/teams/Mythic%20Tech/pipelines/otp/jobs/test-develop/badge)](https://ci.mythic.tech/teams/Mythic%20Tech/pipelines/otp) |
| Build | [![Go Tests](https://ci.mythic.tech/api/v1/teams/Mythic%20Tech/pipelines/otp/jobs/build-develop/badge)](https://ci.mythic.tech/teams/Mythic%20Tech/pipelines/otp) || -->
