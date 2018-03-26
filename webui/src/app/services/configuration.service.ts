import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';

@Injectable()
export class ConfigurationService {

  public ServerWithUrl: string;
  public ServerWithApiUrl: string;
  public AuthToken: string;

  private Server: string = environment.baseUrlUi + "/";
  private ApiUrl: string = environment.baseUrlUi + "/api/";

  /**
   * constructor
   */
  constructor() {
    this.ServerWithUrl = this.Server;
    this.ServerWithApiUrl = this.ApiUrl;
  }

  /**
   * fix session token
   * @param AuthToken 
   */
  public setAuthToken(AuthToken: string): void {
    this.AuthToken = AuthToken;
  }

  /**
   * get token
   * @param AuthToken 
   */
  public getAuthToken(): string {
    return this.AuthToken;
  }

}
