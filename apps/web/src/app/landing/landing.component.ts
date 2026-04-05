import { Component } from '@angular/core';
import { NavbarComponent } from './components/navbar.component';
import { HeroComponent } from './components/hero.component';
import { FeaturesComponent } from './components/features.component';
import { FooterComponent } from './components/footer.component';
import { TrustedByComponent } from './components/trusted-by.component';

@Component({
  selector: 'app-landing',
  standalone: true,
  imports: [NavbarComponent, HeroComponent, FeaturesComponent, FooterComponent, TrustedByComponent],
  templateUrl: './landing.component.html',
})
export class LandingComponent {}
