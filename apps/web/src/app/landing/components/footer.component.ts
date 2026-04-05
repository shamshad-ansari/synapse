import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { LucideAngularModule, Twitter, Linkedin, Github } from 'lucide-angular';

@Component({
  selector: 'app-footer',
  standalone: true,
  imports: [RouterLink, LucideAngularModule],
  templateUrl: './footer.component.html',
})
export class FooterComponent {
  readonly Twitter = Twitter;
  readonly Linkedin = Linkedin;
  readonly Github = Github;
  readonly logoImage = '/assets/synapse-logo.png';
  readonly currentYear = new Date().getFullYear();
}
